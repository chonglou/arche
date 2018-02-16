package nut

import "github.com/astaxie/beego/orm"

// AdminIndexUsers list users
// @router /admin/users [get]
func (p *API) AdminIndexUsers() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var items []User
		if _, err := orm.NewOrm().QueryTable(new(User)).
			OrderBy("-current_sign_in_at").
			All(&items,
				"id", "email", "name",
				"provider_type",
				"sign_in_count",
				"last_sign_in_ip", "last_sign_in_at",
				"current_sign_in_ip", "current_sign_in_at",
			); err != nil {
			return nil, err
		}
		return items, nil
	})
}
