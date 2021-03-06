package nut

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/mux"
	"github.com/chonglou/arche/web/queue"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

func (p *Plugin) deleteUsersSignOut(c *mux.Context) {
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	l := c.Locale()
	if err := p.Dao.AddLog(p.DB, user.ID, c.ClientIP(), l, "nut.logs.user.sign-out"); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) getUsersProfile(c *mux.Context) {
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{
		"name":  user.Name,
		"email": user.Email,
		"logo":  user.Logo,
	})
}

type fmUserProfile struct {
	Name string `json:"name" binding:"required"`
	Logo string `json:"logo" binding:"required"`
}

func (p *Plugin) postUsersProfile(c *mux.Context) {
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmUserProfile
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	user.Name = fm.Name
	user.Logo = fm.Logo
	user.UpdatedAt = time.Now()
	if _, err := p.DB.Model(user).
		Column("name", "logo", "updated_at").
		Update(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

type fmUserChangePassword struct {
	CurrentPassword      string `json:"currentPassword" binding:"required"`
	NewPassword          string `json:"newPassword" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=NewPassword"`
}

func (p *Plugin) postUsersChangePassword(c *mux.Context) {
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmUserChangePassword
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	ip := c.ClientIP()
	l := c.Locale()
	if !p.Security.Check(user.Password, []byte(fm.CurrentPassword)) {
		c.Abort(http.StatusBadRequest, p.I18n.E(l, "nut.errors.user.email-password-not-match"))
		return
	}
	pwd, err := p.Security.Hash([]byte(fm.NewPassword))
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	user.Password = pwd
	user.UpdatedAt = time.Now()
	if err := p.DB.RunInTransaction(func(db *pg.Tx) error {
		if _, err := db.Model(user).Column("password", "updated_at").Update(); err != nil {
			return err
		}
		if err := p.Dao.AddLog(db, user.ID, ip, l, "nut.logs.user.change-password"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) getUsersLogs(c *mux.Context) {
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var items []Log
	if err := p.DB.Model(&items).
		Where("user_id = ?", user.ID).
		Order("created_at DESC").Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type fmUserSignIn struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

func (p *Plugin) postUsersSignIn(c *mux.Context) {
	var fm fmUserSignIn
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	cm := make(jws.Claims)
	l := c.Locale()
	if err := p.DB.RunInTransaction(func(db *pg.Tx) error {
		user, err := p.Dao.SignIn(db, l, c.ClientIP(), fm.Email, fm.Password)
		if err != nil {
			return err
		}
		cm.Set("uid", user.UID)
		cm.Set(RoleAdmin, p.Dao.Is(p.DB, user.ID, RoleAdmin))
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}

	// cm.Set("roles", roles)
	tkn, err := p.Jwt.Sum(cm, time.Hour*24)
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{"token": string(tkn)})
}

type fmUserSignUp struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"email"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *Plugin) postUsersSignUp(c *mux.Context) {
	var fm fmUserSignUp
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}

	ip := c.ClientIP()
	l := c.Locale()
	if err := p.DB.RunInTransaction(func(db *pg.Tx) error {
		user, err := p.Dao.AddEmailUser(db, l, ip, fm.Name, fm.Email, fm.Password)
		if err != nil {
			return err
		}
		if err = p.sendEmail(c, l, user, actConfirm); err != nil {
			log.Error(err)
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, mux.H{})
}

type fmUserEmail struct {
	Email string `json:"email" binding:"email"`
}

func (p *Plugin) postUsersConfirm(c *mux.Context) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	l := c.Locale()
	if user.IsConfirm() {
		c.Abort(http.StatusForbidden, p.I18n.E(l, "nut.errors.user-already-confirm"))
		return
	}
	if err := p.sendEmail(c, l, user, actConfirm); err != nil {
		log.Error(err)
	}

	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) postUsersUnlock(c *mux.Context) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	l := c.Locale()
	if !user.IsLock() {
		c.Abort(http.StatusInternalServerError, p.I18n.E(l, "nut.errors.user.not-lock"))
		return
	}
	if err := p.sendEmail(c, l, user, actUnlock); err != nil {
		log.Error(err)
	}

	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) postUsersForgotPassword(c *mux.Context) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	l := c.Locale()
	if err := p.sendEmail(c, l, user, actResetPassword); err != nil {
		log.Error(err)
	}

	c.JSON(http.StatusOK, mux.H{})
}

type fmUserResetPassword struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *Plugin) postUsersResetPassword(c *mux.Context) {
	var fm fmUserResetPassword
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	cm, err := p.Jwt.Validate([]byte(fm.Token))
	if err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	l := c.Locale()
	if cm.Get("act").(string) != actResetPassword {
		c.Abort(http.StatusInternalServerError, p.I18n.E(l, "errors.bad-action"))
		return
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	pwd, err := p.Security.Hash([]byte(fm.Password))
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	user.Password = pwd
	user.UpdatedAt = time.Now()
	if err = p.DB.RunInTransaction(func(db *pg.Tx) error {
		if _, err = db.Model(user).Column("password", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(db, user.ID, c.ClientIP(), l, "nut.logs.user.reset-password"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) getUsersConfirmToken(c *mux.Context) {
	cm, err := p.Jwt.Validate([]byte(c.Param("token")))
	if err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	l := c.Locale()
	if cm.Get("act").(string) != actConfirm {
		c.Abort(http.StatusBadRequest, p.I18n.E(l, "errors.bad-action"))
		return
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	if user.IsConfirm() {
		c.Abort(http.StatusInternalServerError, p.I18n.E(l, "nut.errors.user-already-confirm"))
		return
	}

	now := time.Now()
	user.UpdatedAt = now
	user.ConfirmedAt = &now
	if err = p.DB.RunInTransaction(func(db *pg.Tx) error {
		if _, err = db.Model(user).Column("confirmed_at", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(db, user.ID, c.ClientIP(), l, "nut.logs.user.confirm"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) getUsersUnlockToken(c *mux.Context) {
	cm, err := p.Jwt.Validate([]byte(c.Param("token")))
	if err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	l := c.Locale()
	if cm.Get("act").(string) != actUnlock {
		c.Abort(http.StatusBadRequest, p.I18n.E(l, "errors.bad-action"))
		return
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	if !user.IsLock() {
		c.Abort(http.StatusInternalServerError, p.I18n.E(l, "nut.errors.user.not-lock"))
		return
	}
	user.LockedAt = nil
	user.UpdatedAt = time.Now()
	if err = p.DB.RunInTransaction(func(db *pg.Tx) error {
		if _, err = db.Model(user).Column("locked_at", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(db, user.ID, c.ClientIP(), l, "nut.logs.unlock"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

const (
	actConfirm       = "confirm"
	actUnlock        = "unlock"
	actResetPassword = "reset-password"

	// SendEmailJob send email
	SendEmailJob = "send.email"
)

func (p *Plugin) sendEmail(c *mux.Context, lang string, user *User, act string) error {
	cm := jws.Claims{}
	cm.Set("act", act)
	cm.Set("uid", user.UID)
	tkn, err := p.Jwt.Sum(cm, time.Hour*6)
	if err != nil {
		return err
	}

	obj := mux.H{
		"home":  c.Home(),
		"token": string(tkn),
	}

	subject, err := p.I18n.H(lang, fmt.Sprintf("nut.users.%s.email-subject", act), obj)
	if err != nil {
		return err
	}
	body, err := p.I18n.H(lang, fmt.Sprintf("nut.users.%s.email-body", act), obj)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(map[string]string{
		"to":      user.Email,
		"subject": subject,
		"body":    body,
	})
	if err != nil {
		return err
	}
	return p.Queue.Put(queue.NewTask(SendEmailJob, 0, buf))
}

func (p *Plugin) doSendEmail(id string, payload []byte) error {
	arg := make(map[string]string)

	if err := json.Unmarshal(payload, &arg); err != nil {
		return err
	}

	to := arg["to"]
	subject := arg["subject"]
	body := arg["body"]
	if viper.GetString("env") != web.PRODUCTION {
		log.Debugf("send to %s: %s\n%s", to, subject, body)
		return nil
	}

	smtp := make(map[string]interface{})
	if err := p.Settings.Get(p.DB, "site.smtp", &smtp); err != nil {
		return err
	}

	sender := smtp["username"].(string)
	msg := gomail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dia := gomail.NewDialer(
		smtp["host"].(string),
		smtp["port"].(int),
		sender,
		smtp["password"].(string),
	)

	return dia.DialAndSend(msg)

}
