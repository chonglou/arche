package nut

import (
	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/mux"
	"github.com/chonglou/arche/web/settings"
	"github.com/go-pg/pg"
)

// HTMLHandlerFunc html handler func
type HTMLHandlerFunc func(string, *mux.Context) error

// RedirectHandlerFunc redirect handle func
type RedirectHandlerFunc func(string, *mux.Context) error

// ObjectHandlerFunc object handle func
type ObjectHandlerFunc func(string, *mux.Context) (interface{}, error)

// Layout layout
type Layout struct {
	I18n     *i18n.I18n         `inject:""`
	Settings *settings.Settings `inject:""`
	Jwt      *web.Jwt           `inject:""`
	DB       *pg.DB             `inject:""`
	Dao      *Dao               `inject:""`
}

// IsAdmin is admin
func (p *Layout) IsAdmin(c *mux.Context) (*User, error) {
	l := c.Locale()
	it, err := p.CurrentUser(c)
	if err != nil {
		return nil, err
	}
	if !p.Dao.Is(p.DB, it.ID, RoleAdmin) {
		return nil, p.I18n.E(l, "errors.not-allow")
	}
	return it, nil
}

// CurrentUser parse user from request
func (p *Layout) CurrentUser(c *mux.Context) (*User, error) {
	l := c.Locale()
	cm, err := p.Jwt.Parse(c.Request)
	if err != nil {
		return nil, err
	}
	uid, ok := cm.Get("uid").(string)
	if !ok {
		return nil, p.I18n.E(l, "errors.not-found")
	}
	it, err := p.Dao.GetUserByUID(p.DB, uid)
	if err != nil {
		return nil, err
	}
	if !it.IsConfirm() || it.IsLock() {
		return nil, p.I18n.E(l, "errors.not-allow")
	}
	return it, nil
}
