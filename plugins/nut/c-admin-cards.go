package nut

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) indexAdminCards(c *mux.Context) {
	l := c.Get(mux.LOCALE).(string)
	var items []Card
	if err := p.DB.Model(&items).Column("id", "loc", "sort", "title", "href").
		Where("lang = ?", l).
		Order("loc ASC").Order("sort ASC").
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type fmCard struct {
	Href    string `json:"href" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Summary string `json:"summary" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Action  string `json:"action" binding:"required"`
	Logo    string `json:"logo" binding:"required"`
	Loc     string `json:"loc" binding:"required"`
	Sort    int    `json:"sort"`
}

func (p *Plugin) createAdminCard(c *mux.Context) {
	l := c.Get(mux.LOCALE).(string)
	var fm fmCard
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	it := Card{
		Href:      fm.Href,
		Title:     fm.Title,
		Summary:   fm.Summary,
		Type:      fm.Type,
		Action:    fm.Action,
		Logo:      fm.Logo,
		Loc:       fm.Loc,
		Sort:      fm.Sort,
		Lang:      l,
		UpdatedAt: time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) showAdminCard(c *mux.Context) {
	var it = Card{}
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) updateAdminCard(c *mux.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	var fm fmCard
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	it := Card{
		ID:        uint(id),
		Title:     fm.Title,
		Type:      fm.Type,
		Action:    fm.Action,
		Summary:   fm.Summary,
		Logo:      fm.Logo,
		Href:      fm.Href,
		Loc:       fm.Loc,
		Sort:      fm.Sort,
		UpdatedAt: time.Now(),
	}

	if _, err := p.DB.Model(&it).Column("title",
		"type",
		"action",
		"summary",
		"logo",
		"href",
		"loc",
		"sort",
		"updated_at",
	).Update(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) destroyAdminCard(c *mux.Context) {
	if _, err := p.DB.Model(new(Card)).
		Where("id = ?", c.Param("id")).
		Delete(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}
