package nut

import "github.com/astaxie/beego"

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
		// TODO favicon
		// TODO current user
		return rst, nil
	})
}
