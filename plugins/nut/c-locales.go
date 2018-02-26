package nut

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func (p *Plugin) indexLocale(_ string, c *gin.Context) (interface{}, error) {
	lng, err := language.Parse(c.Param("lang"))
	if err != nil {
		return nil, err
	}
	items, err := p.I18n.All(lng.String())
	if err != nil {
		return nil, err
	}
	return items, nil
}
