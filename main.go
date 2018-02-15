package main

import (
	"log"
	"os"

	_ "github.com/astaxie/beego/cache/redis"
	"github.com/chonglou/arche/plugins/nut"
	_ "github.com/chonglou/arche/routers"
	_ "github.com/lib/pq"
)

func main() {
	if err := nut.Main(os.Args...); err != nil {
		log.Fatal(err)
	}
}
