package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) getLocales(c *mux.Context) {
	items, err := p.I18n.All(c.Param("lang"))
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}
