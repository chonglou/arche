package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) indexAdminUsers(c *mux.Context) {
	var items []User
	if err := p.DB.Model(&items).Column(
		"id", "email", "name", "provider_type", "logo",
		"sign_in_count", "last_sign_in_at", "last_sign_in_ip", "current_sign_in_at", "current_sign_in_ip",
	).
		Order("last_sign_in_at ASC").
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}
