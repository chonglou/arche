package nut

import (
	"time"

	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/i18n"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Dao dao
type Dao struct {
	Security *web.Security `inject:""`
	I18n     *i18n.I18n    `inject:""`
}

// SignIn set sign-in info
func (p *Dao) SignIn(db orm.DB, lang, ip, email, password string) (*User, error) {
	var it User
	if err := db.Model(&it).
		Where("provider_id = ?", email).
		Where("provider_type = ?", UserTypeEmail).
		Select(); err != nil {
		return nil, err
	}

	if !p.Security.Check(it.Password, []byte(password)) {
		p.AddLog(db, it.ID, ip, lang, "nut.logs.user.sign-in.failed")
		return nil, p.I18n.E(lang, "nut.errors.email-password-not-match")
	}

	if !it.IsConfirm() {
		return nil, p.I18n.E(lang, "nut.errors.user-not-confirm")
	}

	if it.IsLock() {
		return nil, p.I18n.E(lang, "nut.errors.user.is-lock")
	}

	p.AddLog(db, it.ID, ip, lang, "nut.logs.user.sign-in.success")
	it.SignInCount++
	it.LastSignInAt = it.CurrentSignInAt
	it.LastSignInIP = it.CurrentSignInIP
	now := time.Now()
	it.CurrentSignInAt = &now
	it.CurrentSignInIP = ip

	if _, err := db.Model(&it).
		Column(
			"last_sign_in_at", "last_sign_in_ip",
			"current_sign_in_at", "current_sign_in_ip",
			"sign_in_count",
			"updated_at",
		).
		Update(); err != nil {
		return nil, err
	}
	return &it, nil
}

// AddLog add log
func (p *Dao) AddLog(db orm.DB, user uint, ip, lang, format string, args ...interface{}) error {
	err := db.Insert(&Log{
		UserID:  user,
		IP:      ip,
		Message: p.I18n.T(lang, format, args...),
	})
	return err
}

// AddEmailUser add email user
func (p *Dao) AddEmailUser(db orm.DB, lang, ip, name, email, password string) (*User, error) {
	cnt, err := db.Model(new(User)).Where("email = ?", email).Count()
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, p.I18n.E(lang, "nut.errors.email-already-exists")
	}
	passwd, err := p.Security.Hash([]byte(password))
	if err != nil {
		return nil, err
	}
	user := User{
		Email:           email,
		Password:        passwd,
		Name:            name,
		ProviderType:    UserTypeEmail,
		ProviderID:      email,
		LastSignInIP:    "0.0.0.0",
		CurrentSignInIP: "0.0.0.0",
	}
	user.SetUID()
	user.SetGravatarLogo()

	if err := db.Insert(&user); err != nil {
		return nil, err
	}
	if err := p.AddLog(db, user.ID, ip, lang, "nut.logs.user.sign-up"); err != nil {
		return nil, err
	}
	return &user, nil
}

//Is is role ?
func (p *Dao) Is(db orm.DB, user uint, names ...string) bool {
	for _, name := range names {
		if p.Can(db, user, name, DefaultResourceType, DefaultResourceID) {
			return true
		}
	}
	return false
}

//Can can?
func (p *Dao) Can(db orm.DB, user uint, name string, rty string, rid uint) bool {
	var r Role

	if err := db.Model(&r).
		Column("id").
		Where("name = ?", name).
		Where("resource_type = ?", rty).
		Where("resource_id = ?", rid).
		Select(); err != nil {
		return false
	}
	var pm Policy
	if err := db.Model(&pm).
		Column("nbf", "exp").
		Where("user_id = ?", user).
		Where("role_id = ?", r.ID).
		Select(); err != nil {
		return false
	}

	return pm.Enable()
}

// GetRole create role if not exist
func (p *Dao) GetRole(db orm.DB, name string, rty string, rid uint) (*Role, error) {
	it := Role{}
	err := db.Model(&it).
		Where("name = ?", name).
		Where("resource_type = ?", rty).
		Where("resource_id = ?", rid).
		Select()
	if err == nil {
		return &it, nil
	}

	if err != pg.ErrNoRows {
		return nil, err
	}
	it.Name = name
	it.ResourceID = rid
	it.ResourceType = rty
	it.UpdatedAt = time.Now()
	if err = db.Insert(&it); err != nil {
		return nil, err
	}
	return &it, nil
}

//Deny deny permission
func (p *Dao) Deny(db orm.DB, user uint, role uint) error {
	_, err := db.Model(new(Policy)).
		Where("role_id = ?", role).
		Where("user_id = ?", user).
		Delete()
	return err
}

// Authority get roles
func (p *Dao) Authority(db orm.DB, user uint, rty string, rid uint) ([]string, error) {
	var items []*Role

	if err := db.Model(&items).
		Column("id").
		Where("resource_type = ?", rty).
		Where("resource_id = ?", rid).
		Select(); err != nil {
		return nil, err
	}
	var roles []string
	for _, r := range items {
		var pm Policy
		if err := db.Model(&pm).Column("nbf", "exp").
			Where("role_id = ?", r.ID).
			Where("user_id = ?", user).
			Select(); err == nil {
			if pm.Enable() {
				roles = append(roles, r.Name)
			}
		}
	}
	return roles, nil
}

//Allow allow permission
func (p *Dao) Allow(db orm.DB, user uint, role uint, years, months, days int) error {
	now := time.Now()
	exp := now.AddDate(years, months, days)

	var pm Policy
	err := db.Model(&pm).Column("id", "nbf", "exp").
		Where("role_id = ?", role).
		Where("user_id = ?", user).Select()
	pm.Nbf = now
	pm.Exp = exp
	pm.UpdatedAt = now
	if err == nil {
		_, err = db.Model(&pm).Column("nbf", "exp", "updated_at").Update()
		return err
	}
	if err != pg.ErrNoRows {
		return err
	}
	pm.UserID = user
	pm.RoleID = role
	err = db.Insert(&pm)
	return err
}

func (p *Dao) confirmUser(db orm.DB, lang, ip string, user *User) error {
	now := time.Now()
	user.ConfirmedAt = &now
	user.UpdatedAt = now
	if _, err := db.Model(user).Column("confirmed_at", "updated_at").
		Update(); err != nil {
		return err
	}
	return p.AddLog(db, user.ID, ip, lang, "nut.logs.user.confirm")
}

func (p *Dao) setUserPassword(db orm.DB, lang, ip string, user *User, password string) error {
	passwd, err := p.Security.Hash([]byte(password))
	if err != nil {
		return err
	}
	user.Password = passwd
	user.UpdatedAt = time.Now()
	if _, err = db.Model(user).Column("password", "updated_at").Update(); err != nil {
		return err
	}
	return p.AddLog(db, user.ID, ip, lang, "nut.logs.user.change-password")
}
