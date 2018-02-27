package nut

import (
	"github.com/chonglou/arche/web/i18n"
	"github.com/gin-gonic/gin"
)

func (p *Plugin) indexAdminLocales(l string, c *gin.Context) (interface{}, error) {
	var items []i18n.Model
	if err := p.DB.Model(&items).Column("id", "code", "message").
		Where("lang = ?", l).
		Order("code ASC").
		Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmLocale struct {
	Code    string `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (p *Plugin) createAdminLocale(l string, c *gin.Context) (interface{}, error) {
	var fm fmLocale
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	if err := p.I18n.Set(p.DB, l, fm.Code, fm.Message); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) showAdminLocale(l string, c *gin.Context) (interface{}, error) {
	var it i18n.Model
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) destroyAdminLocale(l string, c *gin.Context) (interface{}, error) {
	if _, err := p.DB.Model(new(i18n.Model)).
		Where("id = ?", c.Param("id")).Delete(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
