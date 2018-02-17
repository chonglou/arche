package forum

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/chonglou/arche/plugins/nut"
)

// IndexCatalogs index
// @router /catalogs [get]
func (p *HTML) IndexCatalogs() {
	p.HTML("forum/catalogs/index.html", func() error {
		var items []Catalog
		if _, err := orm.NewOrm().
			QueryTable(new(Catalog)).
			OrderBy("-updated_at").
			All(&items); err != nil {
			return err
		}
		p.Data["catalogs"] = items
		return nil
	})
}

// ShowCatalog show
// @router /catalogs/:id [get]
func (p *HTML) ShowCatalog() {
	p.HTML("forum/catalogs/show.html", func() error {
		var it Catalog
		if err := orm.NewOrm().
			QueryTable(new(Catalog)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return err
		}
		p.Data["catalog"] = it
		return nil
	})
}

// -----------------------------------------------------------------------------

// IndexCatalogs index
// @router /catalogs [get]
func (p *API) IndexCatalogs() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var items []Catalog
		if _, err := orm.NewOrm().
			QueryTable(new(Catalog)).
			OrderBy("-updated_at").
			All(&items, "id", "title", "summary", "color", "icon", "updated_at"); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmCatalog struct {
	Title   string `json:"title" valid:"Required"`
	Summary string `json:"summary" valid:"Required"`
	Color   string `json:"color" valid:"Required"`
	Icon    string `json:"icon" valid:"Required"`
}

// CreateCatalog create
// @router /catalogs [post]
func (p *API) CreateCatalog() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmCatalog
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if _, err := orm.NewOrm().Insert(&Catalog{
			Title:   fm.Title,
			Summary: fm.Summary,
			Icon:    fm.Icon,
			Color:   fm.Color,
		}); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// ShowCatalog show
// @router /catalogs/:id [get]
func (p *API) ShowCatalog() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it Catalog
		if err := orm.NewOrm().
			QueryTable(new(Catalog)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		return it, nil
	})
}

// UpdateCatalog update
// @router /catalogs/:id [post]
func (p *API) UpdateCatalog() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmCatalog
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		var it Catalog
		o := orm.NewOrm()
		if err := o.
			QueryTable(new(Catalog)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		it.Title = fm.Title
		it.Summary = fm.Summary
		it.Color = fm.Color
		it.Icon = fm.Icon
		it.UpdatedAt = time.Now()
		if _, err := o.Update(&it, "title", "summary", "icon", "color", "updated_at"); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// DestroyCatalog destroy
// @router /catalogs/:id [delete]
func (p *API) DestroyCatalog() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		if _, err := orm.NewOrm().QueryTable(new(Catalog)).
			Filter("id", p.Ctx.Input.Param(":id")).
			Delete(); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}
