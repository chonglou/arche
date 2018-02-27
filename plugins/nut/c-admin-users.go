package nut

import "github.com/gin-gonic/gin"

func (p *Plugin) indexAdminUsers(l string, c *gin.Context) (interface{}, error) {
	var items []User
	if err := p.DB.Model(&items).Column(
		"id", "email", "name", "provider_type", "logo",
		"sign_in_count", "last_sign_in_at", "last_sign_in_ip", "current_sign_in_at", "current_sign_in_ip",
	).
		Order("last_sign_in_at ASC").
		Select(); err != nil {
		return nil, err
	}
	return items, nil
}
