package nut

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (p *Plugin) indexAdminFriendLinks(l string, c *gin.Context) (interface{}, error) {
	var items []FriendLink
	if err := p.DB.Model(&items).Order("sort ASC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmFriendLink struct {
	Title string `json:"title" binding:"required"`
	Home  string `json:"home" binding:"required"`
	Logo  string `json:"logo" binding:"required"`
	Sort  int    `json:"sort"`
}

func (p *Plugin) createAdminFriendLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmFriendLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := FriendLink{
		Title:     fm.Title,
		Home:      fm.Home,
		Logo:      fm.Logo,
		Sort:      fm.Sort,
		UpdatedAt: time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showAdminFriendLink(l string, c *gin.Context) (interface{}, error) {
	var it = FriendLink{}
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).Select(); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateAdminFriendLink(l string, c *gin.Context) (interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}
	var fm fmFriendLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
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
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyAdminFriendLink(l string, c *gin.Context) (interface{}, error) {
	if _, err := p.DB.Model(new(FriendLink)).
		Where("id = ?", c.Param("id")).
		Delete(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
