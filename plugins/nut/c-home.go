package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/mux"
	"github.com/spf13/viper"
)

func (p *Plugin) getLayout(c *mux.Context) {
	l := c.Get(mux.LOCALE).(string)
	// site info
	site := mux.H{}
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
	site["locale"] = l
	site["languages"] = viper.GetStringSlice("languages")

	// current-user
	user, ok := c.Get(CurrentUser).(*User)
	// nav
	if ok {
		site["user"] = mux.H{
			"name":  user.Name,
			"type":  user.ProviderType,
			"logo":  user.Logo,
			"admin": p.Dao.Is(p.DB, user.ID, RoleAdmin),
		}
	}

	c.JSON(http.StatusOK, site)
}
