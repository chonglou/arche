package nut

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) indexAdminFriendLinks(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var items []FriendLink
	if err := p.DB.Model(&items).Order("sort ASC").Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type fmFriendLink struct {
	Title string `json:"title" binding:"required"`
	Home  string `json:"home" binding:"required"`
	Logo  string `json:"logo" binding:"required"`
	Sort  int    `json:"sort"`
}

func (p *Plugin) createAdminFriendLink(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmFriendLink
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	it := FriendLink{
		Title:     fm.Title,
		Home:      fm.Home,
		Logo:      fm.Logo,
		Sort:      fm.Sort,
		UpdatedAt: time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) showAdminFriendLink(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var it = FriendLink{}
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) updateAdminFriendLink(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	var fm fmFriendLink
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	it := FriendLink{
		ID:        uint(id),
		Title:     fm.Title,
		Home:      fm.Home,
		Logo:      fm.Logo,
		Sort:      fm.Sort,
		UpdatedAt: time.Now(),
	}

	if _, err := p.DB.Model(&it).
		Column("title", "home", "logo", "sort", "updated_at").
		Update(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) destroyAdminFriendLink(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	if _, err := p.DB.Model(new(FriendLink)).
		Where("id = ?", c.Param("id")).
		Delete(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}
