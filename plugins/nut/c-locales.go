package nut

import "github.com/astaxie/beego/orm"

// IndexLocales list all locales
// @router /locales/:lang [get]
func (p *API) IndexLocales() {
	p.JSON(func() (interface{}, error) {
		items, err := AllLocales(orm.NewOrm(), p.Ctx.Input.Param(":lang"))
		return items, err
	})
}
