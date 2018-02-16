package nut

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// Home home
// @router / [get]
func (p *HTML) Home() {
	p.TplName = "nut/index.html"
}

// GetLayout layout
// @router /layout [get]
func (p *API) GetLayout() {
	p.JSON(func() (interface{}, error) {
		rst := H{
			"locale":    p.Lang,
			"languages": beego.AppConfig.Strings("languages"),
		}
		// site info
		for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
			rst[k] = Tr(p.Lang, "site."+k)
		}

		// site author
		o := orm.NewOrm()
		author := make(map[string]string)
		if err := Get(o, "site.author", &author); err != nil {
			author["email"] = "manager@change-me.com"
			author["name"] = "who-am-i"
		}
		rst["author"] = author
		// favicon
		var favicon string
		if err := Get(o, "site.favicon", &favicon); err == nil {
			rst["favicon"] = favicon
		} else {
			rst["favicon"] = "/assets/favicon.png"
		}

		return rst, nil
	})
}
