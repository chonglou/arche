package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/mux"
	"github.com/go-pg/pg"
)

func (p *Plugin) postInstall(c *mux.Context) {
	var fm fmUserSignUp
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}

	ip := c.ClientIP()
	l := c.Get(mux.LOCALE).(string)
	if err := p.DB.RunInTransaction(func(db *pg.Tx) error {
		cnt, err := db.Model(new(User)).Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			return p.I18n.E(l, "errors.forbidden")
		}
		user, err := p.Dao.AddEmailUser(db, l, ip, fm.Name, fm.Email, fm.Password)
		if err != nil {
			return err
		}
		if err = p.Dao.confirmUser(db, l, ip, user); err != nil {
			return err
		}
		for _, r := range []string{RoleAdmin, RoleRoot} {
			role, err := p.Dao.GetRole(db, r, DefaultResourceType, DefaultResourceID)
			if err != nil {
				return err
			}
			if err = p.Dao.Allow(db, user.ID, role.ID, 50, 0, 0); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, mux.H{})
}
