package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) getLayout(c *mux.Context) {
	lang := c.Get(i18n.LOCALE).(string)
	// site info
	site := mux.H{}
	for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
		site[k] = p.I18n.T(lang, "site."+k)
	}
	author := make(map[string]string)
	p.Settings.Get(p.DB, "site.author", &author)
	site["author"] = author

	// favicon
	var favicon string
	p.Settings.Get(p.DB, "site.favicon", &favicon)
	site["favicon"] = favicon

	// i18n
	site[i18n.LOCALE] = lang
	site["languages"] = p.I18n.Languages()

	// current-user
	user := c.Get(CurrentUser)
	// nav
	if user != nil {
		user := user.(*User)
		site["user"] = mux.H{
			"name":  user.Name,
			"type":  user.ProviderType,
			"admin": c.Get(IsAdmin).(bool),
		}
	}

	c.JSON(http.StatusOK, site)
}
