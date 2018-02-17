package forum

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/chonglou/arche/plugins/nut"
)

// IndexTags index
// @router /tags [get]
func (p *HTML) IndexTags() {
	p.HTML("forum/tags/index.html", func() error {
		var items []Tag
		if _, err := orm.NewOrm().
			QueryTable(new(Tag)).
			OrderBy("-updated_at").
			All(&items); err != nil {
			return err
		}
		p.Data["tags"] = items
		return nil
	})
}

// ShowTag show
// @router /tags/:id [get]
func (p *HTML) ShowTag() {
	p.HTML("forum/tags/show.html", func() error {
		var it Tag
		if err := orm.NewOrm().
			QueryTable(new(Tag)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return err
		}
		p.Data["tag"] = it
		return nil
	})
}

// -----------------------------------------------------------------------------

// IndexTags index
// @router /tags [get]
func (p *API) IndexTags() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var items []Tag
		if _, err := orm.NewOrm().
			QueryTable(new(Tag)).
			OrderBy("-updated_at").
			All(&items, "id", "name", "color", "updated_at"); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmTag struct {
	Name  string `json:"name" valid:"Required"`
	Color string `json:"color" valid:"Required"`
}

// CreateTag create
// @router /tags [post]
func (p *API) CreateTag() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmTag
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if _, err := orm.NewOrm().Insert(&Tag{
			Name:  fm.Name,
			Color: fm.Color,
		}); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// ShowTag show
// @router /tags/:id [get]
func (p *API) ShowTag() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it Tag
		if err := orm.NewOrm().
			QueryTable(new(Tag)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		return it, nil
	})
}

// UpdateTag update
// @router /tags/:id [post]
func (p *API) UpdateTag() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmTag
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		var it Tag
		o := orm.NewOrm()
		if err := o.
			QueryTable(new(Tag)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		it.Name = fm.Name
		it.Color = fm.Color
		it.UpdatedAt = time.Now()
		if _, err := o.Update(&it, "name", "color", "updated_at"); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// DestroyTag destroy
// @router /tags/:id [delete]
func (p *API) DestroyTag() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		if _, err := orm.NewOrm().QueryTable(new(Tag)).
			Filter("id", p.Ctx.Input.Param(":id")).
			Delete(); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}
