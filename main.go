package main

import (
	"log"
	"os"

	_ "github.com/chonglou/arche/plugins/blog"
	_ "github.com/chonglou/arche/plugins/calendar"
	_ "github.com/chonglou/arche/plugins/cbeta"
	_ "github.com/chonglou/arche/plugins/dict"
	_ "github.com/chonglou/arche/plugins/donate"
	_ "github.com/chonglou/arche/plugins/forum"
	_ "github.com/chonglou/arche/plugins/hotel"
	_ "github.com/chonglou/arche/plugins/library"
	_ "github.com/chonglou/arche/plugins/nut"
	_ "github.com/chonglou/arche/plugins/ops/mail"
	_ "github.com/chonglou/arche/plugins/ops/vpn"
	_ "github.com/chonglou/arche/plugins/shop"
	_ "github.com/chonglou/arche/plugins/survey"
	_ "github.com/chonglou/arche/plugins/wiki"
	"github.com/chonglou/arche/web"
)

func main() {
	if err := web.Main(os.Args...); err != nil {
		log.Fatal(err)
	}
}
