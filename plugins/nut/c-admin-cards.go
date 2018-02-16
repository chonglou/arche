package nut

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// AdminIndexCards index
// @router /admin/cards [get]
func (p *API) AdminIndexCards() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var items []Card
		if _, err := orm.NewOrm().QueryTable(new(Card)).
			Filter("lang", p.Lang).
			OrderBy("loc", "sort").
			All(&items); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmCard struct {
	Loc     string `json:"loc" valid:"Required"`
	Sort    int    `json:"sort"`
	Title   string `json:"title" valid:"Required"`
	Logo    string `json:"logo" valid:"Required"`
	Summary string `json:"summary" valid:"Required"`
	Type    string `json:"type" valid:"Required"`
	Action  string `json:"action" valid:"Required"`
	Href    string `json:"href" valid:"Required"`
}

// AdminCreateCard create
// @router /admin/cards [post]
func (p *API) AdminCreateCard() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmCard
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		if _, err := orm.NewOrm().Insert(&Card{
			Loc:     fm.Loc,
			Sort:    fm.Sort,
			Title:   fm.Title,
			Logo:    fm.Logo,
			Summary: fm.Summary,
			Type:    fm.Type,
			Action:  fm.Action,
			Href:    fm.Href,
			Lang:    p.Lang,
		}); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

// AdminShowCard show
// @router /admin/cards/:id [get]
func (p *API) AdminShowCard() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it Card
		if err := orm.NewOrm().QueryTable(new(Card)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		return it, nil
	})
}

// AdminUpdateCard update
// @router /admin/cards/:id [post]
func (p *API) AdminUpdateCard() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmCard
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		o := orm.NewOrm()
		var it Card
		if err := o.QueryTable(new(Card)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		it.Loc = fm.Loc
		it.Sort = fm.Sort
		it.Title = fm.Title
		it.Logo = fm.Logo
		it.Summary = fm.Summary
		it.Type = fm.Type
		it.Action = fm.Action
		it.Href = fm.Href
		it.UpdatedAt = time.Now()
		if _, err := o.Update(&it,
			"loc", "sort",
			"title", "logo", "summary", "type", "action", "href",
			"updated_at",
		); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

// AdminDestroyCard destroy
// @router /admin/cards/:id [delete]
func (p *API) AdminDestroyCard() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it Card
		if _, err := orm.NewOrm().QueryTable(new(Card)).
			Filter("id", p.Ctx.Input.Param(":id")).
			Delete(); err != nil {
			return nil, err
		}
		return it, nil
	})
}
