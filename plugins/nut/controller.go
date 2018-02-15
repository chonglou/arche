package nut

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
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
	p.setLocale()
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

// BindJson bind to json data
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
