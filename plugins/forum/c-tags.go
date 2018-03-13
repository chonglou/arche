package forum

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

// -----------------------------------------------------------------------------

func (p *Plugin) indexTags(l string, c *gin.Context) (interface{}, error) {
	var items []Tag
	if err := p.DB.Model(&items).
		Column("id", "name", "color").
		Order("updated_at DESC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmTag struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required"`
}

func (p *Plugin) createTag(l string, c *gin.Context) (interface{}, error) {
	var fm fmTag
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Tag{
		Name:      fm.Name,
		Color:     fm.Color,
		UpdatedAt: time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showTag(l string, c *gin.Context) (interface{}, error) {
	var it = Tag{}
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateTag(l string, c *gin.Context) (interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}
	var fm fmTag
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Tag{
		ID:        uint(id),
		Name:      fm.Name,
		Color:     fm.Color,
		UpdatedAt: time.Now(),
	}

	if _, err := p.DB.Model(&it).
		Column("name", "color", "updated_at").
		Update(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyTag(l string, c *gin.Context) (interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}
	if er := p.DB.RunInTransaction(func(db *pg.Tx) error {
		if _, e := db.Model(new(Tag)).
			Column("article.id", "Articles").
			Where("id = ?", id).
			Delete(); e != nil {
			return e
		}
		return nil
	}); er != nil {
		return nil, er
	}
	return gin.H{}, err
}
