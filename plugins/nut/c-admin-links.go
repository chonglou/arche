package nut

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) indexAdminLinks(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	l := c.Locale()
	var items []Link
	if err := p.DB.Model(&items).Column("id", "label", "href", "loc", "x", "y").
		Where("lang = ?", l).
		Order("loc ASC").Order("x ASC").Order("y ASC").
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type fmLink struct {
	Href  string `json:"href" binding:"required"`
	Label string `json:"label" binding:"required"`
	Loc   string `json:"loc" binding:"required"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

func (p *Plugin) createAdminLink(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmLink
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	l := c.Locale()
	it := Link{
		Href:      fm.Href,
		Label:     fm.Label,
		Loc:       fm.Loc,
		X:         fm.X,
		Y:         fm.Y,
		Lang:      l,
		UpdatedAt: time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) showAdminLink(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var it = Link{}
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) updateAdminLink(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}

	var fm fmLink
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	it := Link{
		ID:        uint(id),
		Label:     fm.Label,
		Href:      fm.Href,
		Loc:       fm.Loc,
		X:         fm.X,
		Y:         fm.Y,
		UpdatedAt: time.Now(),
	}

	if _, err := p.DB.Model(&it).
		Column("label", "href", "loc", "x", "y", "updated_at").
		Update(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) destroyAdminLink(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	if _, err := p.DB.Model(new(Link)).
		Where("id = ?", c.Param("id")).
		Delete(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}
