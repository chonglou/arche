package nut

import (
	"strings"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/SermoDigital/jose/jws"
	"github.com/astaxie/beego/orm"
)

// DeleteUsersSignOut sign out
// @router /users/sign-out [delete]
func (p *API) DeleteUsersSignOut() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		if err := AddLog(orm.NewOrm(),
			user.ID, p.Ctx.Input.IP(),
			p.Lang, "nut.logs.user.sign-out"); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

type fmUserChangePassword struct {
	CurrentPassword      string `json:"currentPassword" valid:"Required"`
	NewPassword          string `json:"newPassword" valid:"Required;MinSize(6);MaxSize(32)"`
	PasswordConfirmation string `json:"passwordConfirmation" valid:"Required"`
}

// PostUsersChangePassword change user password
// @router /users/change-password [post]
func (p *API) PostUsersChangePassword() {
	p.JSON(func() (interface{}, error) {
		var fm fmUserChangePassword
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if fm.NewPassword != fm.PasswordConfirmation {
			return nil, Te(p.Lang, "nut.errors.passwords-not-match")
		}

		user := p.MustSignIn()
		if !user.Auth(user.Email, fm.CurrentPassword) {
			return nil, Te(p.Lang, "nut.errors.email-password-not-match")
		}
		if err := user.SetPassword(fm.NewPassword); err != nil {
			return nil, err
		}
		user.UpdatedAt = time.Now()
		var items []Log
		o := orm.NewOrm()
		o.Begin()
		if _, err := o.Update(user, "password", "updated_at"); err != nil {
			o.Rollback()
			return nil, err
		}
		if err := AddLog(o, user.ID, p.Ctx.Input.IP(), p.Lang, "nut.logs.user.change-password"); err != nil {
			return nil, err
		}
		o.Commit()
		return items, nil
	})
}

type fmUserProfile struct {
	Logo string `json:"logo" valid:"Required"`
	Name string `json:"name	" valid:"Required"`
}

// PostUsersProfile set user profile
// @router /users/profile [post]
func (p *API) PostUsersProfile() {
	p.JSON(func() (interface{}, error) {
		var fm fmUserProfile
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		user := p.MustSignIn()
		user.Logo = fm.Logo
		user.Name = fm.Name
		user.UpdatedAt = time.Now()
		var items []Log
		if _, err := orm.NewOrm().
			Update(user, "logo", "name", "updated_at"); err != nil {
			return nil, err
		}
		return items, nil
	})
}

// GetUsersProfile get user profile
// @router /users/profile [get]
func (p *API) GetUsersProfile() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		return H{
			"email": user.Email,
			"name":  user.Name,
			"logo":  user.Logo,
		}, nil
	})
}

// GetUsersLogs get user logs
// @router /users/logs [get]
func (p *API) GetUsersLogs() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		var items []Log
		if _, err := orm.NewOrm().QueryTable(new(Log)).
			Filter("user_id", user.ID).
			OrderBy("-created_at").
			All(&items); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmUserSignIn struct {
	Email    string `json:"email" valid:"Email"`
	Password string `json:"password" valid:"Required"`
}

// PostUsersSignIn sign in
// @router /users/sign-in [post]
func (p *API) PostUsersSignIn() {
	p.JSON(func() (interface{}, error) {
		var fm fmUserSignIn
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		var it User
		o := orm.NewOrm()
		if err := o.QueryTable(new(User)).
			Filter("email", strings.ToLower(fm.Email)).
			One(&it); err != nil {
			return nil, err
		}
		ip := p.Ctx.Input.IP()
		o.Begin()
		if !it.Auth(fm.Email, fm.Password) {
			if err := AddLog(o, it.ID, ip, p.Lang, "nut.logs.user.sign-in.failed"); err != nil {
				o.Rollback()
				return nil, err
			}
			return nil, Te(p.Lang, "nut.errors.email-password-not-match")
		}

		if !it.IsConfirm() {
			return nil, Te(p.Lang, "nut.errors.user-not-confirm")
		}

		if it.IsLock() {
			return nil, Te(p.Lang, "nut.errors.user-is-lock")
		}
		if err := signIn(o, p.Lang, ip, &it); err != nil {
			o.Rollback()
			return nil, err
		}
		o.Commit()

		cm := make(jws.Claims)
		cm.Set("uid", it.UID)
		cm.Set(RoleAdmin, Is(o, it.ID, RoleAdmin))
		tkn, err := JWT().Sum(cm, time.Hour*24)
		if err != nil {
			return nil, err
		}
		return H{"token": string(tkn)}, nil
	})
}

type fmUserSignUp struct {
	Name                 string `json:"name" valid:"Required"`
	Email                string `json:"email" valid:"Email;MaxSize(255)"`
	Password             string `json:"password" valid:"Required;MinSize(6);MaxSize(32)"`
	PasswordConfirmation string `json:"passwordConfirmation" valid:"Required"`
}

// PostUsersSignUp sign up
// @router /users/sign-up [post]
func (p *API) PostUsersSignUp() {
	p.JSON(func() (interface{}, error) {
		o := orm.NewOrm()
		var fm fmUserSignUp
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if fm.Password != fm.PasswordConfirmation {
			return nil, Te(p.Lang, "nut.errors.passwords-not-match")
		}

		o.Begin()
		user, err := addEmailUser(o, p.Ctx.Input.IP(), p.Lang, fm.Name, fm.Email, fm.Password)
		if err != nil {
			o.Rollback()
			return nil, err
		}
		o.Commit()
		if err := p.sendEmail(user, actConfirm); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

type fmUserEmail struct {
	Email string `json:"email" valid:"Email;MaxSize(255)"`
}

// PostUsersConfirm confirm
// @router /users/confirm [post]
func (p *API) PostUsersConfirm() {
	p.JSON(func() (interface{}, error) {
		var fm fmUserEmail
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		var it User
		if err := orm.NewOrm().QueryTable(new(User)).
			Filter("email", fm.Email).
			One(&it, "email", "uid", "confirmed_at"); err != nil {
			return nil, err
		}

		if it.IsConfirm() {
			return nil, Te(p.Lang, "nut.errors.user-already-confirm")
		}
		if err := p.sendEmail(&it, actConfirm); err != nil {
			return nil, err
		}

		return H{}, nil
	})
}

// PostUsersUnlock unlock
// @router /users/unlock [post]
func (p *API) PostUsersUnlock() {
	p.JSON(func() (interface{}, error) {
		var fm fmUserEmail
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		var it User
		if err := orm.NewOrm().QueryTable(new(User)).
			Filter("email", fm.Email).
			One(&it, "email", "uid", "locked_at"); err != nil {
			return nil, err
		}

		if !it.IsLock() {
			return nil, Te(p.Lang, "nut.errors.user-not-lock")
		}
		if err := p.sendEmail(&it, actUnlock); err != nil {
			return nil, err
		}

		return H{}, nil
	})
}

// PostUsersForgotPassword forgot password
// @router /users/forgot-password [post]
func (p *API) PostUsersForgotPassword() {
	p.JSON(func() (interface{}, error) {
		var fm fmUserEmail
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		var it User
		if err := orm.NewOrm().QueryTable(new(User)).
			Filter("email", fm.Email).
			One(&it, "email", "uid"); err != nil {
			return nil, err
		}

		if err := p.sendEmail(&it, actResetPassword); err != nil {
			return nil, err
		}

		return H{}, nil

	})
}

// GetUsersConfirmToken confirm
// @router /users/confirm/:token [get]
func (p *HTML) GetUsersConfirmToken() {
	p.Redirect("/", func() error {
		cm, err := JWT().Validate([]byte(p.Ctx.Input.Param(":token")))
		if err != nil {
			return err
		}
		if cm.Get("act").(string) != actConfirm {
			return Te(p.Lang, "errors.forbidden")
		}

		o := orm.NewOrm()
		var it User
		err = o.QueryTable(new(User)).
			Filter("uid", cm.Get("uid").(string)).
			One(&it, "id", "confirmed_at")
		if err != nil {
			return err
		}
		if it.IsConfirm() {
			return Te(p.Lang, "nut.errors.user-already-confirm")
		}

		o.Begin()
		if err = confirmUser(o, p.Ctx.Input.IP(), p.Lang, it.ID); err != nil {
			o.Rollback()
			return err
		}
		o.Commit()

		return nil
	})
}

// GetUsersUnlockToken unlock
// @router /users/unlock/:token [get]
func (p *HTML) GetUsersUnlockToken() {
	p.Redirect("/", func() error {
		cm, err := JWT().Validate([]byte(p.Ctx.Input.Param(":token")))
		if err != nil {
			return err
		}
		if cm.Get("act").(string) != actUnlock {
			return Te(p.Lang, "errors.forbidden")
		}

		o := orm.NewOrm()
		var it User
		err = o.QueryTable(new(User)).
			Filter("uid", cm.Get("uid").(string)).
			One(&it, "id", "locked_at")
		if err != nil {
			return err
		}
		if !it.IsLock() {
			return Te(p.Lang, "nut.errors.user-not-lock")
		}

		o.Begin()

		it.UpdatedAt = time.Now()
		it.LockedAt = nil
		if _, err = o.Update(&it, "locked_at", "updated_at"); err != nil {
			o.Rollback()
			return err
		}
		if err = AddLog(o, it.ID, p.Ctx.Input.IP(), p.Lang, "nut.logs.user.unlock"); err != nil {
			o.Rollback()
			return err
		}
		o.Commit()

		return nil
	})
}

type fmUserResetPassword struct {
	Token                string `json:"token" valid:"Required"`
	Password             string `json:"password" valid:"Required;MinSize(6);MaxSize(32)"`
	PasswordConfirmation string `json:"passwordConfirmation" valid:"Required"`
}

// PostUsersResetPassword reset password
// @router /users/reset-password [post]
func (p *API) PostUsersResetPassword() {
	p.JSON(func() (interface{}, error) {
		var fm fmUserResetPassword
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if fm.Password != fm.PasswordConfirmation {
			return nil, Te(p.Lang, "nut.errors.passwords-not-match")
		}
		cm, err := JWT().Validate([]byte(fm.Token))
		if err != nil {
			return nil, err
		}
		if cm.Get("act").(string) != actResetPassword {
			return nil, Te(p.Lang, "errors.forbidden")
		}

		o := orm.NewOrm()
		var it User
		err = o.QueryTable(new(User)).
			Filter("uid", cm.Get("uid").(string)).
			One(&it, "id")
		if err != nil {
			return nil, err
		}
		if err = it.SetPassword(fm.Password); err != nil {
			return nil, err
		}
		it.UpdatedAt = time.Now()

		o.Begin()
		if _, err = o.Update(&it, "password", "updated_at"); err != nil {
			o.Rollback()
			return nil, err
		}
		if err = AddLog(o, it.ID, p.Ctx.Input.IP(), p.Lang, "nut.logs.user.reset-password"); err != nil {
			o.Rollback()
			return nil, err
		}
		o.Commit()

		return H{}, nil
	})
}

const (
	actConfirm       = "nut.users.confirm"
	actUnlock        = "nut.users.unlock"
	actResetPassword = "nut.users.reset-password"
)

func (p *API) sendEmail(user *User, act string) error {
	cm := jws.Claims{}
	cm.Set("act", act)
	cm.Set("uid", user.UID)
	tkn, err := JWT().Sum(cm, time.Hour*6)
	if err != nil {
		return err
	}

	obj := H{
		"home":  p.Ctx.Input.Site(),
		"token": string(tkn),
	}

	subject, err := Th(p.Lang, act+".email-subject", obj)
	if err != nil {
		return err
	}
	body, err := Th(p.Lang, act+".email-body", obj)
	if err != nil {
		return err
	}

	_, err = Server().SendTask(&tasks.Signature{
		Name: GetFunctionName(SendEmailTask),
		Args: []tasks.Arg{
			{Type: "string", Value: user.Email},
			{Type: "string", Value: subject},
			{Type: "string", Value: body},
		},
	})
	return err

}
