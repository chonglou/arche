package nut

import (
	"net/http"

	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/mux"
	"github.com/chonglou/arche/web/settings"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

const (
	// CurrentUser current user
	CurrentUser = "currentUser"
	// IsAdmin is admin?
	IsAdmin = "isAdmin"
)

// Layout layout
type Layout struct {
	I18n     *i18n.I18n         `inject:""`
	Settings *settings.Settings `inject:""`
	Jwt      *web.Jwt           `inject:""`
	DB       *pg.DB             `inject:""`
	Dao      *Dao               `inject:""`
}

// MustAdminMiddleware must-admin
func (p *Layout) MustAdminMiddleware(c *mux.Context) {
	if is := c.Get(IsAdmin); is.(bool) {
		return
	}
	l := c.Get(mux.LOCALE).(string)
	c.Abort(http.StatusForbidden, p.I18n.E(l, "errors.not-allow"))
}

// MustSignInMiddleware must-sign-in
func (p *Layout) MustSignInMiddleware(c *mux.Context) {
	if _, ok := c.Get(CurrentUser).(*User); ok {
		return
	}
	l := c.Get(mux.LOCALE).(string)
	c.Abort(http.StatusForbidden, p.I18n.E(l, "errors.not-allow"))
}

// CurrentUserMiddleware parse user from request
func (p *Layout) CurrentUserMiddleware(c *mux.Context) {
	cm, err := p.Jwt.Parse(c.Request)
	if err != nil {
		log.Error(err)
		return
	}
	uid, ok := cm.Get("uid").(string)
	if !ok {
		log.Error("bad token")
		return
	}
	it, err := p.Dao.GetUserByUID(p.DB, uid)
	if err != nil {
		log.Error(err)
		return
	}
	if !it.IsConfirm() || it.IsLock() {
		log.Error("bad user status")
		return
	}
	c.Set(CurrentUser, it)
	c.Set(IsAdmin, p.Dao.Is(p.DB, it.ID, RoleAdmin))
	return
}
