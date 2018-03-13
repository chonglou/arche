package forum

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

// -----------------------------------------------------------------------------
func (p *Plugin) indexCatalogs(l string, c *gin.Context) (interface{}, error) {
	var items []Catalog
	if err := p.DB.Model(&items).
		Column("id", "title", "summary", "color", "icon").
		Order("updated_at DESC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmCatalog struct {
	Title   string `json:"title" validate:"required"`
	Summary string `json:"summary" validate:"required"`
	Icon    string `json:"icon" validate:"required"`
	Color   string `json:"color" validate:"required"`
}

func (p *Plugin) createCatalog(l string, c *gin.Context) (interface{}, error) {
	var fm fmCatalog
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Catalog{
		Title:     fm.Title,
		Summary:   fm.Summary,
		Icon:      fm.Icon,
		Color:     fm.Color,
		UpdatedAt: time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showCatalog(l string, c *gin.Context) (interface{}, error) {
	var it = Catalog{}
	if err := p.DB.Model(&it).
		Where("id = ?", c.Param("id")).
		Select(); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateCatalog(l string, c *gin.Context) (interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}
	var fm fmCatalog
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Catalog{
		ID:        uint(id),
		Title:     fm.Title,
		Summary:   fm.Summary,
		Icon:      fm.Icon,
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

func (p *Plugin) destroyCatalog(l string, c *gin.Context) (interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}
	if err = p.DB.RunInTransaction(func(db *pg.Tx) error {
		if _, er := db.Model(new(Topic)).
			Where("catalog_id = ?", id).
			Delete(); er != nil {
			return er
		}
		if _, er := db.Model(new(Catalog)).
			Where("id = ?", id).
			Delete(); er != nil {
			return er
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return gin.H{}, err
}
