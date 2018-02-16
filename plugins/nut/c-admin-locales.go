package nut

import "github.com/astaxie/beego/orm"

// GetAdminLocales list all locales
// @router /admin/locales [get]
func (p *API) GetAdminLocales() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		items, err := AllLocales(orm.NewOrm(), p.Lang)
		return items, err
	})
}

type fmLocale struct {
	Code    string `json:"code" valid:"Required"`
	Message string `json:"message" valid:"Required"`
}

// PostAdminLocales list all locales
// @router /admin/locales [post]
func (p *API) PostAdminLocales() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmLocale
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		err := SetLocale(orm.NewOrm(), p.Lang, fm.Code, fm.Message)
		return H{}, err
	})
}

// GetAdminLocale get locales
// @router /admin/locales/:code [get]
func (p *API) GetAdminLocale() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		return H{
			"message": Tr(p.Lang, p.Ctx.Input.Param(":code")),
		}, nil
	})
}

// DeleteAdminLocale delete locales
// @router /admin/locales/:code [delete]
func (p *API) DeleteAdminLocale() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		_, err := orm.NewOrm().QueryTable(new(Locale)).
			Filter("code", p.Ctx.Input.Param(":code")).
			Filter("lang", p.Lang).
			Delete()
		return H{}, err
	})
}
