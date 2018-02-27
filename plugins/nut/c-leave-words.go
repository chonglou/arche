package nut

import "github.com/gin-gonic/gin"

func (p *Plugin) createLeaveWord(l string, c *gin.Context) (interface{}, error) {
	var fm fmLeaveWord
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := LeaveWord{
		Body: fm.Body,
		Type: fm.Type,
	}
	if err := p.DB.Insert(&it); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) indexAdminLeaveWords(l string, c *gin.Context) (interface{}, error) {
	var items []LeaveWord
	if err := p.DB.Model(&items).Order("created_at DESC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmLeaveWord struct {
	Body string `json:"body" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (p *Plugin) destroyAdminLeaveWord(l string, c *gin.Context) (interface{}, error) {
	if _, err := p.DB.Model(new(LeaveWord)).
		Where("id = ?", c.Param("id")).
		Delete(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
