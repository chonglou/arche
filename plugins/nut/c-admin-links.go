package nut

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// AdminIndexLinks index
// @router /admin/links [get]
func (p *API) AdminIndexLinks() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var items []Link
		if _, err := orm.NewOrm().QueryTable(new(Link)).
			Filter("lang", p.Lang).
			OrderBy("loc", "x", "y").
			All(&items); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmLink struct {
	Loc   string `json:"loc" valid:"Required"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Label string `json:"label" valid:"Required"`
	Href  string `json:"href" valid:"Required"`
}

// AdminCreateLink create
// @router /admin/links [post]
func (p *API) AdminCreateLink() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmLink
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		if _, err := orm.NewOrm().Insert(&Link{
			Loc:   fm.Loc,
			Label: fm.Label,
			Href:  fm.Href,
			X:     fm.X,
			Y:     fm.Y,
			Lang:  p.Lang,
		}); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

// AdminShowLink show
// @router /admin/links/:id [get]
func (p *API) AdminShowLink() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it Link
		if err := orm.NewOrm().QueryTable(new(Link)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		return it, nil
	})
}

// AdminUpdateLink update
// @router /admin/links/:id [post]
func (p *API) AdminUpdateLink() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmLink
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		o := orm.NewOrm()
		var it Link
		if err := o.QueryTable(new(Link)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		it.Loc = fm.Loc
		it.X = fm.X
		it.Y = fm.Y
		it.Label = fm.Label
		it.Href = fm.Href
		it.UpdatedAt = time.Now()
		if _, err := o.Update(&it,
			"loc", "x", "y",
			"label", "href",
			"updated_at",
		); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

// AdminDestroyLink destroy
// @router /admin/links/:id [delete]
func (p *API) AdminDestroyLink() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it Link
		if _, err := orm.NewOrm().QueryTable(new(Link)).
			Filter("id", p.Ctx.Input.Param(":id")).
			Delete(); err != nil {
			return nil, err
		}
		return it, nil
	})
}
