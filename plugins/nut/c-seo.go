package nut

import (
	"net/http"
	"path"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func (p *Plugin) getSitemapXMLGz(c *gin.Context) {
	if err := p.Sitemap.ToXMLGz(p.Layout.Home(c), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
func (p *Plugin) getRssAtom(c *gin.Context) {
	author := make(map[string]string)
	p.Settings.Get(p.DB, "site.author", &author)
	lang := c.Param("lang")
	if err := p.RSS.ToAtom(
		p.Layout.Home(c), lang,
		p.I18n.T(lang, "site.title"), p.I18n.T(lang, "site.description"),
		&feeds.Author{
			Name:  author["name"],
			Email: author["email"],
		},
		c.Writer,
	); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

	}
}
func (p *Plugin) getRobotsTxt(c *gin.Context) {
	tpl, err := template.New("robots.txt").ParseFiles(path.Join("templates", "robots.txt"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := tpl.Execute(c.Writer, gin.H{"home": p.Layout.Home(c)}); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
