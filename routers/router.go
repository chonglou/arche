package routers

import (
	"github.com/astaxie/beego"
	"github.com/chonglou/arche/plugins/nut"
)

func init() {
	beego.Include(&nut.HTML{})
	beego.AddNamespace(beego.NewNamespace("/api", beego.NSInclude(&nut.API{})))
}
