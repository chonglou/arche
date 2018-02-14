package nut

import "github.com/astaxie/beego/orm"

// IndexLocales list all locales
// @router /locales [get]
func (p *API) IndexLocales() {
	items, err := AllLocales(orm.NewOrm(), p.Lang)
	p.Check(err)
	p.JSON(items)
}
