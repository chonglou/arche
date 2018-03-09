package donate

import (
	"time"

	"github.com/chonglou/arche/plugins/nut"
	"github.com/gin-gonic/gin"
)

func (p *Plugin) getProjects(l string, c *gin.Context) error {
	var items []Project
	if err := p.DB.Model(&items).
		Column("id", "title").
		Order("updated_at DESC").
		Select(); err != nil {
		return err
	}
	c.Set("projects", items)
	return nil
}

func (p *Plugin) getProject(l string, c *gin.Context) error {
	var it Project
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		return err
	}
	c.Set("project", it)
	return nil
}

// -----------------------------------------------------------------------------

func (p *Plugin) indexProjects(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	admin := c.MustGet(nut.IsAdmin).(bool)
	var items []Project
	qry := p.DB.Model(&items).Column("id", "title")
	if !admin {
		qry = qry.Where("user_id", user.ID)
	}
	if err := qry.Order("updated_at DESC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmProject struct {
	Title   string `json:"title" binding:"required"`
	Body    string `json:"body" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Methods string `json:"methods" binding:"required"`
}

func (p *Plugin) createProject(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmProject
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
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
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showProject(l string, c *gin.Context) (interface{}, error) {
	it, err := p.canEditProject(l, c)
	return it, err
}

func (p *Plugin) updateProject(l string, c *gin.Context) (interface{}, error) {
	var fm fmProject
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it, err := p.canEditProject(l, c)
	if err != nil {
		return nil, err
	}

	it.Title = fm.Title
	it.Body = fm.Body
	it.Type = fm.Type
	it.Methods = fm.Methods
	it.UpdatedAt = time.Now()

	if _, err := p.DB.Model(it).
		Column("title", "body", "type", "methods", "updated_at").
		Update(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyProject(l string, c *gin.Context) (interface{}, error) {
	it, err := p.canEditProject(l, c)
	if err != nil {
		return nil, err
	}
	if err := p.DB.Delete(it); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) canEditProject(l string, c *gin.Context) (*Project, error) {
	var it Project
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).
		Select(); err != nil {
		return nil, err
	}
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	admin := c.MustGet(nut.IsAdmin).(bool)
	if it.UserID == user.ID || admin {
		return &it, nil
	}
	return nil, p.I18n.E(l, "errors.not-allow")
}
