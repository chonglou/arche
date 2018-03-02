package nut

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (p *Plugin) indexAdminCards(l string, c *gin.Context) (interface{}, error) {
	var items []Card
	if err := p.DB.Model(&items).Column("id", "loc", "sort", "title", "href").
		Where("lang = ?", l).
		Order("loc ASC").Order("sort ASC").
		Select(); err != nil {
		return nil, err
	}
	return items, nil
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

func (p *Plugin) createAdminCard(l string, c *gin.Context) (interface{}, error) {
	var fm fmCard
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
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
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showAdminCard(l string, c *gin.Context) (interface{}, error) {
	var it = Card{}
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateAdminCard(l string, c *gin.Context) (interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}
	var fm fmCard
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
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
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyAdminCard(l string, c *gin.Context) (interface{}, error) {
	if _, err := p.DB.Model(new(Card)).
		Where("id = ?", c.Param("id")).
		Delete(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
