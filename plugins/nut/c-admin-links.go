package nut

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (p *Plugin) indexAdminLinks(l string, c *gin.Context) (interface{}, error) {
	var items []Link
	if err := p.DB.Model(&items).Column("id", "label", "href", "loc", "x", "y").
		Where("lang = ?", l).
		Order("loc ASC, x ASC, y ASC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmLink struct {
	Href  string `json:"href" binding:"required"`
	Label string `json:"label" binding:"required"`
	Loc   string `json:"loc" binding:"required"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

func (p *Plugin) createAdminLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
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
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showAdminLink(l string, c *gin.Context) (interface{}, error) {
	var it = Link{}
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateAdminLink(l string, c *gin.Context) (interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}

	var fm fmLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
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
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyAdminLink(l string, c *gin.Context) (interface{}, error) {
	if _, err := p.DB.Model(new(Link)).Where("id = ?", c.Param("id")).Delete(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
