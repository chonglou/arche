package routers

import (
	"github.com/astaxie/beego"
	"github.com/chonglou/arche/plugins/donate"
	"github.com/chonglou/arche/plugins/forum"
	"github.com/chonglou/arche/plugins/nut"
)

func init() {

	beego.Include(&nut.HTML{})
	for k, v := range map[string]beego.ControllerInterface{
		"donate": &donate.HTML{},
		"forum":  &forum.HTML{},
	} {
		beego.AddNamespace(beego.NewNamespace("/"+k, beego.NSInclude(v)))
	}

	// --------------------------------

	api := []beego.LinkNamespace{
		beego.NSInclude(&nut.API{}),
	}
	for k, v := range map[string]beego.ControllerInterface{
		"donate": &donate.API{},
		"forum":  &forum.API{},
	} {
		api = append(api, beego.NSNamespace("/"+k, beego.NSInclude(v)))
	}
	beego.AddNamespace(beego.NewNamespace("/api", api...))
}
