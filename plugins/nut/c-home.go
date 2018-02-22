package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) getHome(c *mux.Context) {
	c.JSON(http.StatusOK, mux.H{})
}
