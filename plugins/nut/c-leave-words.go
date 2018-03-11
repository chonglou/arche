package nut

import (
	"net/http"

	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) createLeaveWord(c *mux.Context) {
	var fm fmLeaveWord
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	it := LeaveWord{
		Body: fm.Body,
		Type: fm.Type,
	}
	if err := p.DB.Insert(&it); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) indexAdminLeaveWords(c *mux.Context) {
	var items []LeaveWord
	if err := p.DB.Model(&items).Order("created_at DESC").Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type fmLeaveWord struct {
	Body string `json:"body" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (p *Plugin) destroyAdminLeaveWord(c *mux.Context) {
	if _, err := p.DB.Model(new(LeaveWord)).
		Where("id = ?", c.Param("id")).
		Delete(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}
