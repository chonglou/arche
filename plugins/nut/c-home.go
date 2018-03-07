package nut

import (
	"github.com/chonglou/arche/web/i18n"
	"github.com/gin-gonic/gin"
)

func (p *Plugin) getHome(ctx *gin.Context) {
	wrap := func(f HTMLHandlerFunc) HTMLHandlerFunc {
		return func(l string, d gin.H, c *gin.Context) error {
			// google site verify code
			var googleVerifyCode string
			p.Settings.Get(p.DB, googleSiteVerification, &googleVerifyCode)
			d["googleVerifyCode"] = googleVerifyCode
			if f == nil {
				return nil
			}
			return f(l, d, c)
		}
	}

	var home string
	if err := p.Settings.Get(p.DB, "site.home", &home); err == nil {
		if hnd := p.HomePage.Get(home); hnd != nil {
			p.Layout.HTML(home, wrap(hnd))(ctx)
			return
		}
	}

	p.Layout.HTML("nut/index", wrap(nil))(ctx)
}

func (p *Plugin) getLayout(l string, c *gin.Context) (interface{}, error) {
	// site info
	site := gin.H{}
	for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
		site[k] = p.I18n.T(l, "site."+k)
	}
	author := make(map[string]string)
	p.Settings.Get(p.DB, "site.author", &author)
	site["author"] = author

	// favicon
	var favicon string
	p.Settings.Get(p.DB, "site.favicon", &favicon)
	site["favicon"] = favicon

	// i18n
	site[i18n.LOCALE] = l
	site["languages"], _ = p.I18n.Languages()

	// current-user
	user, ok := c.Get(CurrentUser)
	// nav
	if ok {
		user := user.(*User)
		site["user"] = gin.H{
			"name":  user.Name,
			"type":  user.ProviderType,
			"logo":  user.Logo,
			"admin": c.MustGet(IsAdmin).(bool),
		}
	}

	return site, nil
}
