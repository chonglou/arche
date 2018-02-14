package nut

import (
	"math"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"golang.org/x/text/language"
)

// Controller controller
type Controller struct {
	beego.Controller

	Lang string
}

// Check check error
func (p *Controller) Check(err error) {
	if err != nil {
		p.CustomAbort(http.StatusInternalServerError, err.Error())
	}
}

// Prepare runs after Init before request function execution.
func (p *Controller) Prepare() {
	p.setLocale()
}

// JSON render json
func (p *Controller) JSON(v interface{}) {
	p.Data["json"] = v
	p.ServeJSON()
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
