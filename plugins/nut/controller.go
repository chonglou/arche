package nut

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"
	"golang.org/x/text/language"
)

// Controller controller
type Controller struct {
	beego.Controller

	Lang string
}

// Prepare runs after Init before request function execution.
func (p *Controller) Prepare() {
	p.Layout = "layouts/application/index.html"
	p.setLocale()
}

// MustSignIn must sign in
func (p *Controller) MustSignIn() *User {
	user, err := p.CurrentUser()
	if err != nil {
		p.CustomAbort(http.StatusForbidden, err.Error())
		return nil
	}
	return user
}

// MustAdmin must admin
func (p *Controller) MustAdmin() *User {
	user := p.MustSignIn()
	if !Is(orm.NewOrm(), user.ID, RoleAdmin) {
		p.CustomAbort(http.StatusForbidden, Tr(p.Lang, "errors.forbidden"))
		return nil
	}
	return user
}

// CurrentUser get current user
func (p *Controller) CurrentUser() (*User, error) {
	cm, err := JWT().Parse(p.Ctx.Request)
	if err != nil {
		return nil, err
	}
	uid, ok := cm.Get("uid").(string)
	if !ok {
		return nil, Te(p.Lang, "errors.forbidden")
	}

	var it User
	o := orm.NewOrm()
	if err := o.QueryTable(new(User)).
		Filter("uid", uid).
		One(&it,
			"id", "uid",
			"provider_type", "name", "email",
			"current_sign_in_at", "locked_at", "confirmed_at",
		); err != nil {
		return nil, err
	}
	if !it.IsConfirm() {
		return nil, Te(p.Lang, "nut.errors.user-not-confirm")
	}

	if it.IsLock() {
		return nil, Te(p.Lang, "nut.errors.user-is-lock")
	}

	return &it, nil
}

func (p *Controller) valid(v interface{}) error {
	var valid validation.Validation
	ok, err := valid.Valid(v)
	if ok {
		return nil
	}
	if err != nil {
		return err
	}

	var msg []string
	for _, it := range valid.Errors {
		msg = append(msg, it.Message)
	}
	return errors.New(strings.Join(msg, "\n"))
}

// BindJSON bind to json data
func (p *Controller) BindJSON(v interface{}) error {
	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(v); err != nil {
		return err
	}
	return p.valid(v)
}

// BindForm bind to form data
func (p *Controller) BindForm(v interface{}) error {
	if err := p.ParseForm(v); err != nil {
		return err
	}
	return p.valid(v)
}

// HTML render html
func (p *Controller) HTML(tpl string, f func() error) {
	if err := f(); err == nil {
		p.TplName = tpl
	} else {
		beego.Error(err)
		p.TplName = "nut/error.html"
		p.Data["reason"] = err.Error()
	}
}

// Redirect redirect
func (p *Controller) Redirect(u string, f func() error) {
	if e := f(); e == nil {
		p.Controller.Redirect(u, http.StatusFound)
	} else {
		beego.Error(e)
		p.TplName = "nut/error.html"
		p.Data["reason"] = e.Error()
	}
}

// JSON render json
func (p *Controller) JSON(f func() (interface{}, error)) {
	if v, e := f(); e == nil {
		p.Data["json"] = v
		p.ServeJSON()
	} else {
		beego.Error(e)
		p.CustomAbort(http.StatusInternalServerError, e.Error())
	}
}

func (p *Controller) setLocale() {
	const key = "locale"
	var write bool

	// 1. Check URL arguments.
	lang := p.Input().Get(key)

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = p.Ctx.GetCookie(key)
	} else {
		write = true
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := p.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			lang = al[:5] // Only compare first 5 letters.
		}

		write = true
	}

	// 4. Default language is English.
	if len(lang) == 0 || !i18n.IsExist(lang) {
		lang = language.English.String()
		write = true
	}

	// Save language information in cookies.
	if write {
		beego.Debug(key, lang)
		p.Ctx.SetCookie(key, lang, math.MaxUint32, "/")
	}

	// Set language properties.
	p.Lang = lang
	p.Data[key] = lang
	p.Data["languages"] = beego.AppConfig.Strings("languages")

}
