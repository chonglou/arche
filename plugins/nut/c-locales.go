package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/mux"
	"golang.org/x/text/language"
)

func (p *Plugin) indexLocale(c *mux.Context) {
	lng, err := language.Parse(c.Param("lang"))
	if err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	items, err := p.I18n.All(lng.String())
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}
