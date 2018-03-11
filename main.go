package main

import (
	"log"
	"os"

	_ "github.com/chonglou/arche/plugins/calendar"
	_ "github.com/chonglou/arche/plugins/cbeta"
	_ "github.com/chonglou/arche/plugins/dict"
	_ "github.com/chonglou/arche/plugins/donate"
	// _ "github.com/chonglou/arche/plugins/forum"
	_ "github.com/chonglou/arche/plugins/hotel"
	_ "github.com/chonglou/arche/plugins/library"
	_ "github.com/chonglou/arche/plugins/nut"
	_ "github.com/chonglou/arche/plugins/ops/mail"
	_ "github.com/chonglou/arche/plugins/ops/vpn"
	_ "github.com/chonglou/arche/plugins/shop"
	_ "github.com/chonglou/arche/plugins/survey"
	"github.com/chonglou/arche/web"
	_ "github.com/chonglou/arche/web/hugo/ananke"
	_ "github.com/chonglou/arche/web/hugo/even"
	_ "github.com/chonglou/arche/web/hugo/initio"
	_ "github.com/chonglou/arche/web/hugo/kube"
)

func main() {
	if err := web.Main(os.Args...); err != nil {
		log.Fatal(err)
	}
}
