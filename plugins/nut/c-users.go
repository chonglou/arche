package nut

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/chonglou/arche/web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

func (p *Plugin) deleteUsersSignOut(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	if err := p.Dao.AddLog(p.DB, user.ID, c.ClientIP(), l, "nut.logs.user.sign-out"); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) getUsersProfile(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	return gin.H{
		"name":  user.Name,
		"email": user.Email,
		"logo":  user.Logo,
	}, nil
}

type fmUserProfile struct {
	Name string `json:"name" binding:"required"`
	Logo string `json:"logo" binding:"required"`
}

func (p *Plugin) postUsersProfile(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserProfile
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	user := c.MustGet(CurrentUser).(*User)
	user.Name = fm.Name
	user.Logo = fm.Logo
	user.UpdatedAt = time.Now()
	if _, err := p.DB.Model(user).
		Column("name", "logo", "updated_at").
		Update(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

type fmUserChangePassword struct {
	CurrentPassword      string `json:"currentPassword" binding:"required"`
	NewPassword          string `json:"newPassword" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=NewPassword"`
}

func (p *Plugin) postUsersChangePassword(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserChangePassword
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	user := c.MustGet(CurrentUser).(*User)
	ip := c.ClientIP()
	if !p.Security.Check(user.Password, []byte(fm.CurrentPassword)) {
		return nil, p.I18n.E(l, "nut.errors.user.email-password-not-match")
	}
	pwd, err := p.Security.Hash([]byte(fm.NewPassword))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return gin.H{}, nil
}

func (p *Plugin) getUsersLogs(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	var items []Log
	if err := p.DB.Model(&items).Column("message", "created_at").
		Where("user_id = ?", user.ID).
		Order("created_at DESC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmUserSignIn struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

func (p *Plugin) postUsersSignIn(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserSignIn
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	cm := make(jws.Claims)
	if err := p.DB.RunInTransaction(func(db *pg.Tx) error {
		user, err := p.Dao.SignIn(db, l, c.ClientIP(), fm.Email, fm.Password)
		if err != nil {
			return err
		}
		cm.Set(UID, user.UID)
		cm.Set(RoleAdmin, p.Dao.Is(p.DB, user.ID, RoleAdmin))
		return nil
	}); err != nil {
		return nil, err
	}

	// cm.Set("roles", roles)
	tkn, err := p.Jwt.Sum(cm, time.Hour*24)
	if err != nil {
		return nil, err
	}
	return gin.H{"token": string(tkn)}, nil
}

type fmUserSignUp struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"email"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *Plugin) postUsersSignUp(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserSignUp
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	ip := c.ClientIP()
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
		return nil, err
	}

	return gin.H{}, nil
}

type fmUserEmail struct {
	Email string `json:"email" binding:"email"`
}

func (p *Plugin) postUsersConfirm(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		return nil, err
	}
	if user.IsConfirm() {
		return nil, p.I18n.E(l, "nut.errors.user-already-confirm")
	}
	if err := p.sendEmail(c, l, user, actConfirm); err != nil {
		log.Error(err)
	}

	return gin.H{}, nil
}

func (p *Plugin) postUsersUnlock(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		return nil, err
	}
	if !user.IsLock() {
		return nil, p.I18n.E(l, "nut.errors.user.not-lock")
	}
	if err := p.sendEmail(c, l, user, actUnlock); err != nil {
		log.Error(err)
	}

	return gin.H{}, nil
}

func (p *Plugin) postUsersForgotPassword(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		return nil, err
	}
	if err := p.sendEmail(c, l, user, actResetPassword); err != nil {
		log.Error(err)
	}

	return gin.H{}, nil
}

type fmUserResetPassword struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *Plugin) postUsersResetPassword(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserResetPassword
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	cm, err := p.Jwt.Validate([]byte(fm.Token))
	if err != nil {
		return nil, err
	}
	if cm.Get("act").(string) != actResetPassword {
		return nil, p.I18n.E(l, "errors.bad-action")
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		return nil, err
	}
	pwd, err := p.Security.Hash([]byte(fm.Password))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return gin.H{}, nil
}

func (p *Plugin) getUsersConfirmToken(l string, c *gin.Context) error {
	cm, err := p.Jwt.Validate([]byte(c.Param("token")))
	if err != nil {
		return err
	}
	if cm.Get("act").(string) != actConfirm {
		return p.I18n.E(l, "errors.bad-action")
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		return err
	}
	if user.IsConfirm() {
		return p.I18n.E(l, "nut.errors.user-already-confirm")
	}

	now := time.Now()
	user.UpdatedAt = now
	user.ConfirmedAt = &now
	err = p.DB.RunInTransaction(func(db *pg.Tx) error {
		if _, err = db.Model(user).Column("confirmed_at", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(db, user.ID, c.ClientIP(), l, "nut.logs.user.confirm"); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		ss := sessions.Default(c)
		ss.AddFlash(p.I18n.T(l, "flash.success"), NOTICE)
		ss.Save()
	}
	return err
}

func (p *Plugin) getUsersUnlockToken(l string, c *gin.Context) error {
	cm, err := p.Jwt.Validate([]byte(c.Param("token")))
	if err != nil {
		return err
	}
	if cm.Get("act").(string) != actUnlock {
		return p.I18n.E(l, "errors.bad-action")
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		return err
	}
	if !user.IsLock() {
		return p.I18n.E(l, "nut.errors.user.not-lock")
	}
	user.LockedAt = nil
	user.UpdatedAt = time.Now()
	err = p.DB.RunInTransaction(func(db *pg.Tx) error {
		if _, err = db.Model(user).Column("locked_at", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(db, user.ID, c.ClientIP(), l, "nut.logs.unlock"); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		ss := sessions.Default(c)
		ss.AddFlash(p.I18n.T(l, "flash.success"), NOTICE)
		ss.Save()
	}
	return err
}

const (
	actConfirm       = "confirm"
	actUnlock        = "unlock"
	actResetPassword = "reset-password"

	// SendEmailJob send email
	SendEmailJob = "send.email"
)

func (p *Plugin) sendEmail(c *gin.Context, lang string, user *User, act string) error {
	cm := jws.Claims{}
	cm.Set("act", act)
	cm.Set("uid", user.UID)
	tkn, err := p.Jwt.Sum(cm, time.Hour*6)
	if err != nil {
		return err
	}

	obj := gin.H{
		"home":  p.Layout.Home(c),
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
	return p.Queue.Put(SendEmailJob, uuid.New().String(), 0, buf)
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
