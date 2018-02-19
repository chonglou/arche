package main

import (
	"log"
	"os"

	"github.com/chonglou/arche/plugins/nut"
)

func main() {
	if err := nut.Main(os.Args...); err != nil {
		log.Fatal(err)
	}
}
