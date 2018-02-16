package nut

import (
	"encoding/json"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/orm"
)

const (
	// RoleAdmin admin role
	RoleAdmin = "admin"
	// RoleRoot root role
	RoleRoot = "root"

	// UserTypeEmail email type
	UserTypeEmail = "email"

	// DefaultResourceType default resource type
	DefaultResourceType = ""
	// DefaultResourceID default resourc id
	DefaultResourceID = 0
)

// Set settings set
func Set(o orm.Ormer, k string, v interface{}, f bool) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if f {
		buf, err = AES().Encrypt(buf)
		if err != nil {
			return err
		}
	}
	var it Setting
	err = o.QueryTable(new(Setting)).Filter("key", k).One(&it, "id")
	it.Encode = f
	it.Value = string(buf)

	if err == nil {
		_, err = o.Update(&it, "encode", "value")
	} else if err == orm.ErrNoRows {
		it.Key = k
		_, err = o.Insert(&it)
	}
	return err
}

// Get settings get
func Get(o orm.Ormer, k string, v interface{}) error {
	var it Setting
	err := o.QueryTable(new(Setting)).Filter("key", k).One(&it, "id", "value", "encode")
	if err != nil {
		return err
	}
	buf := []byte(it.Value)
	if it.Encode {
		if buf, err = AES().Decrypt(buf); err != nil {
			return err
		}
	}
	return json.Unmarshal(buf, v)
}

//Is is role ?
func Is(o orm.Ormer, user uint, names ...string) bool {
	for _, name := range names {
		if Can(o, user, name, DefaultResourceType, DefaultResourceID) {
			return true
		}
	}
	return false
}

//Can can?
func Can(o orm.Ormer, user uint, name string, rty string, rid uint) bool {
	role, err := getRole(o, name, rty, rid)
	if err != nil {
		return false
	}
	var pm Policy
	if err := o.QueryTable(new(Policy)).
		Filter("user_id", user).
		Filter("role_id", role.ID).
		One(&pm, "nbf", "exp"); err != nil {
		return false
	}

	return pm.Enable()
}

func getRole(o orm.Ormer, name string, rty string, rid uint) (*Role, error) {
	var it Role
	err := o.QueryTable(new(Role)).
		Filter("name", name).
		Filter("resource_type", rty).
		Filter("resource_id", rid).One(&it, "id")
	if err == nil {
		return &it, nil
	}
	if err == orm.ErrNoRows {
		it.Name = name
		it.ResourceType = rty
		it.ResourceID = rid
		if _, err = o.Insert(&it); err == nil {
			return &it, err
		}
	}
	return nil, err
}

// Deny deny permission
func Deny(o orm.Ormer, ip, lang string, user uint, name string, rty string, rid uint) error {
	role, err := getRole(o, name, rty, rid)
	if err != nil {
		return err
	}
	if _, err = o.QueryTable(new(Policy)).
		Filter("user_id", user).
		Filter("role_id", role.ID).
		Delete(); err != nil {
		return err
	}
	return AddLog(o, user, ip, lang, "nut.logs.user.deny", role)
}

// Authority get roles
func Authority(o orm.Ormer, user uint, rty string, rid uint) ([]string, error) {
	var policies []*Policy
	if _, err := o.QueryTable(new(Policy)).
		Filter("resource_type", rty, rid).
		Filter("resource_id", rid).
		Filter("user_id", user).
		All(&policies, "role_id", "nbf", "exp"); err != nil {
		return nil, err
	}
	var roles []string
	for _, it := range policies {
		if it.Enable() {
			if err := o.QueryTable(new(Role)).
				Filter("id", it.Role.ID).
				One(it.Role, "name"); err != nil {
				return nil, err
			}
			roles = append(roles, it.Role.Name)
		}
	}
	return roles, nil
}

//Apply apply permission
func Apply(o orm.Ormer, ip, lang string, user uint, name string, rty string, rid uint, years, months, days int) error {
	now := time.Now()
	exp := now.AddDate(years, months, days)

	role, err := getRole(o, name, rty, rid)
	if err != nil {
		return err
	}

	var it Policy
	err = o.QueryTable(new(Policy)).
		Filter("role_id", role.ID).
		Filter("user_id", user).One(&it, "id")
	if err == nil {
		it.Nbf = now
		it.Exp = exp
		it.UpdatedAt = now
		_, err = o.Update(&it, "nbf", "exp", "updated_at")
	} else if err == orm.ErrNoRows {
		it.User = &User{ID: user}
		it.Role = role
		it.Nbf = now
		it.Exp = exp
		_, err = o.Insert(&it)
	}
	if err != nil {
		return err
	}
	return AddLog(o, user, ip, lang, "nut.logs.user.apply", role)
}

func signIn(o orm.Ormer, lang, ip string, user *User) error {
	if err := AddLog(o, user.ID, ip, lang, "nut.logs.user.sign-in.success"); err != nil {
		return err
	}
	user.SignInCount++
	user.LastSignInAt = user.CurrentSignInAt
	user.LastSignInIP = user.CurrentSignInIP
	now := time.Now()
	user.CurrentSignInAt = &now
	user.CurrentSignInIP = ip
	user.UpdatedAt = now

	_, err := o.Update(user,
		"last_sign_in_at",
		"last_sign_in_ip",
		"current_sign_in_at",
		"current_sign_in_ip",
		"sign_in_count",
		"updated_at",
	)
	return err
}

func addEmailUser(o orm.Ormer, ip, lang, name, email, password string) (*User, error) {
	user := User{
		Name:            name,
		ProviderType:    UserTypeEmail,
		ProviderID:      email,
		LastSignInIP:    "0.0.0.0",
		CurrentSignInIP: "0.0.0.0",
	}
	user.SetEmail(email)
	user.SetUID()
	user.SetGravatarLogo()
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	if cnt, err := o.QueryTable(new(User)).Filter("email", email).Count(); err != nil {
		return nil, err
	} else if cnt > 0 {
		return nil, Te(lang, "errors.email-already-exists")
	}

	if _, err := o.Insert(&user); err != nil {
		return nil, err
	}
	if err := AddLog(o, user.ID, ip, lang, "nut.logs.user.sign-up"); err != nil {
		return nil, err
	}
	return &user, nil
}

// AddLog add log
func AddLog(o orm.Ormer, user uint, ip, lang, code string, args ...interface{}) error {
	_, err := o.Insert(&Log{
		User:    &User{ID: user},
		IP:      ip,
		Message: Tr(lang, code, args...),
	})
	return err
}

func confirmUser(o orm.Ormer, ip, lang string, user uint) error {
	now := time.Now()
	if _, err := o.QueryTable(new(User)).Filter("id", user).Update(
		orm.Params{
			"confirmed_at": now,
			"updated_at":   now,
		},
	); err != nil {
		return err
	}
	return AddLog(o, user, ip, lang, "nut.logs.user.confirm")
}

// AllLocales locale map by lang
func AllLocales(o orm.Ormer, lang string) (map[string]string, error) {
	rst := make(map[string]string)
	// load from file
	cfg, err := goconfig.LoadConfigFile(localeFile(lang))
	if err != nil {
		return nil, err
	}
	for _, s := range cfg.GetSectionList() {
		it, er := cfg.GetSection(s)
		if er != nil {
			return nil, er
		}
		for k, v := range it {
			rst[s+"."+k] = v
		}
	}
	// load from database
	var items []Locale
	if _, err := o.QueryTable(new(Locale)).
		Filter("lang", lang).
		All(&items, "code", "message"); err != nil {
		return nil, err
	}
	for _, it := range items {
		rst[it.Code] = it.Message
	}

	return rst, nil
}

// SetLocale set locale
func SetLocale(o orm.Ormer, lang, code, message string) error {
	var it Locale
	err := o.QueryTable(new(Locale)).
		Filter("lang", lang).
		Filter("code", code).
		One(&it, "id")
	if err == nil {
		it.Message = message
		it.UpdatedAt = time.Now()
		_, err = o.Update(&it, "Message", "UpdatedAt")
	} else if err == orm.ErrNoRows {
		_, err = o.Insert(&Locale{
			Lang:    lang,
			Code:    code,
			Message: message,
		})
	}
	return err
}

// GetLocale get locale message
func GetLocale(o orm.Ormer, lang, code string) (string, error) {
	var it Locale
	if err := o.QueryTable(new(Locale)).
		Filter("lang", lang).
		Filter("code", code).
		One(&it, "message"); err != nil {
		return "", err
	}
	return it.Message, nil
}
