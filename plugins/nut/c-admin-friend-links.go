package nut

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// AdminIndexFriendLinks index
// @router /admin/friend-links [get]
func (p *API) AdminIndexFriendLinks() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var items []FriendLink
		if _, err := orm.NewOrm().QueryTable(new(FriendLink)).
			OrderBy("sort").
			All(&items); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmFriendLink struct {
	Title string `json:"title" valid:"Required"`
	Logo  string `json:"logo" valid:"Required"`
	Home  string `json:"home" valid:"Required"`
	Sort  int    `json:"sort"`
}

// AdminCreateFriendLink create
// @router /admin/friend-links [post]
func (p *API) AdminCreateFriendLink() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmFriendLink
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		if _, err := orm.NewOrm().Insert(&FriendLink{
			Title: fm.Title,
			Logo:  fm.Logo,
			Home:  fm.Home,
			Sort:  fm.Sort,
		}); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

// AdminShowFriendLink show
// @router /admin/friend-links/:id [get]
func (p *API) AdminShowFriendLink() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it FriendLink
		if err := orm.NewOrm().QueryTable(new(FriendLink)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		return it, nil
	})
}

// AdminUpdateFriendLink update
// @router /admin/friend-links/:id [post]
func (p *API) AdminUpdateFriendLink() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmFriendLink
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		o := orm.NewOrm()
		var it FriendLink
		if err := o.QueryTable(new(FriendLink)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return nil, err
		}
		it.Title = fm.Title
		it.Logo = fm.Logo
		it.Home = fm.Home
		it.Sort = fm.Sort
		it.UpdatedAt = time.Now()
		if _, err := o.Update(&it,
			"title", "logo", "home", "sort",
			"updated_at",
		); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

// AdminDestroyFriendLink destroy
// @router /admin/friend-links/:id [delete]
func (p *API) AdminDestroyFriendLink() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var it FriendLink
		if _, err := orm.NewOrm().QueryTable(new(FriendLink)).
			Filter("id", p.Ctx.Input.Param(":id")).
			Delete(); err != nil {
			return nil, err
		}
		return it, nil
	})
}
