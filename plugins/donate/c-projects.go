package donate

import (
	"net/http"
	"time"

	"github.com/chonglou/arche/plugins/nut"
	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) indexProjects(c *mux.Context) {
	var items []Project

	if err := p.DB.Model(&items).
		Column("id", "title").
		Order("updated_at DESC").
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type fmProject struct {
	Title   string `json:"title" validate:"required"`
	Body    string `json:"body" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Methods string `json:"methods" validate:"required"`
}

func (p *Plugin) createProject(c *mux.Context) {
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	var fm fmProject
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	it := Project{
		Title:     fm.Title,
		Body:      fm.Body,
		Type:      fm.Type,
		Methods:   fm.Methods,
		UserID:    user.ID,
		UpdatedAt: time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) showProject(c *mux.Context) {
	var it Project
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) updateProject(c *mux.Context) {
	var fm fmProject
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	it, err := p.canEditProject(c)
	if err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}

	it.Title = fm.Title
	it.Body = fm.Body
	it.Type = fm.Type
	it.Methods = fm.Methods
	it.UpdatedAt = time.Now()

	if _, err := p.DB.Model(it).
		Column("title", "body", "type", "methods", "updated_at").
		Update(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) destroyProject(c *mux.Context) {
	it, err := p.canEditProject(c)
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	if err := p.DB.Delete(it); err != nil {
		c.Abort(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) canEditProject(c *mux.Context) (*Project, error) {
	var it Project
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).
		Select(); err != nil {
		return nil, err
	}
	l := c.Locale()
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		return nil, err
	}

	if it.UserID != user.ID && p.Dao.Is(p.DB, user.ID, nut.RoleAdmin) {
		return nil, p.I18n.E(l, "errors.forbidden")
	}
	return &it, nil
}
