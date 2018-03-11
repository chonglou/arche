package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) indexAdminLocales(c *mux.Context) {
	l := c.Get(mux.LOCALE).(string)
	var items []i18n.Model
	if err := p.DB.Model(&items).Column("id", "code", "message").
		Where("lang = ?", l).
		Order("code ASC").
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type fmLocale struct {
	Code    string `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (p *Plugin) createAdminLocale(c *mux.Context) {
	var fm fmLocale
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	l := c.Get(mux.LOCALE).(string)
	if err := p.I18n.Set(p.DB, l, fm.Code, fm.Message); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) showAdminLocale(c *mux.Context) {
	var it i18n.Model
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) destroyAdminLocale(c *mux.Context) {
	if _, err := p.DB.Model(new(i18n.Model)).
		Where("id = ?", c.Param("id")).Delete(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}
