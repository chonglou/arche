package donate

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/chonglou/arche/plugins/nut"
)

// https://developer.paypal.com/docs/classic/paypal-payments-standard/integration-guide/html_example_donate/

// IndexProjects index
// @router /projects [get]
func (p *HTML) IndexProjects() {
	p.HTML("donate/projects/index.html", func() error {
		var items []Project
		if _, err := orm.NewOrm().
			QueryTable(new(Project)).
			OrderBy("-updated_at").
			All(&items); err != nil {
			return err
		}
		p.Data["projects"] = items
		return nil
	})
}

// ShowProject show
// @router /projects/:id [get]
func (p *HTML) ShowProject() {
	p.HTML("donate/projects/show.html", func() error {
		var it Project
		if err := orm.NewOrm().
			QueryTable(new(Project)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return err
		}
		p.Data["project"] = it
		return nil
	})
}

// -----------------------------------------------------------------------------

// IndexProjects index
// @router /projects [get]
func (p *API) IndexProjects() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var items []Project
		if _, err := orm.NewOrm().
			QueryTable(new(Project)).
			OrderBy("-updated_at").
			All(&items, "id", "title", "updated_at"); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmProject struct {
	Title   string `json:"title" valid:"Required"`
	Body    string `json:"body" valid:"Required"`
	Type    string `json:"type" valid:"Required"`
	Methods string `json:"methods" valid:"Required"`
}

// CreateProject create
// @router /projects [post]
func (p *API) CreateProject() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmProject
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if _, err := orm.NewOrm().Insert(&Project{
			Title:   fm.Title,
			Body:    fm.Body,
			Type:    fm.Type,
			Methods: fm.Methods,
		}); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// ShowProject show
// @router /projects/:id [get]
func (p *API) ShowProject() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it Project
		if err := orm.NewOrm().
			QueryTable(new(Project)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		return it, nil
	})
}

// UpdateProject update
// @router /projects/:id [post]
func (p *API) UpdateProject() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmProject
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		var it Project
		o := orm.NewOrm()
		if err := o.
			QueryTable(new(Project)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		it.Title = fm.Title
		it.Body = fm.Body
		it.Type = fm.Type
		it.Methods = fm.Methods
		it.UpdatedAt = time.Now()
		if _, err := o.Update(&it, "title", "body", "type", "methods", "updated_at"); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// DestroyProject destroy
// @router /projects/:id [delete]
func (p *API) DestroyProject() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		if _, err := orm.NewOrm().QueryTable(new(Project)).
			Filter("id", p.Ctx.Input.Param(":id")).
			Delete(); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}
