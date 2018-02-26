package nut

import (
	"github.com/chonglou/arche/web/i18n"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

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
	langs, err := p.I18n.Languages()
	if err != nil {
		log.Error(err)
		langs = []string{
			language.AmericanEnglish.String(),
			language.SimplifiedChinese.String(),
			language.TraditionalChinese.String(),
		}
	}
	site["languages"] = langs

	// current-user
	user, ok := c.Get(CurrentUser)
	// nav
	if ok {
		user := user.(*User)
		site["user"] = gin.H{
			"name":  user.Name,
			"type":  user.ProviderType,
			"admin": c.MustGet(IsAdmin).(bool),
		}
	}

	return site, nil
}
