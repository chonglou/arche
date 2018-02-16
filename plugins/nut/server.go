package nut

import (
	"fmt"
	"path"
	"sync"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/astaxie/beego"
)

type serverLogger struct {
}

func (p *serverLogger) Print(args ...interface{}) {
	beego.Info(fmt.Sprint(args...))
}
func (p *serverLogger) Printf(f string, args ...interface{}) {
	beego.Info(fmt.Sprintf(f, args...))
}
func (p *serverLogger) Println(args ...interface{}) {
	beego.Info(fmt.Sprintln(args...))
}

func (p *serverLogger) Fatal(args ...interface{}) {
	beego.Error(fmt.Sprint(args...))
}
func (p *serverLogger) Fatalf(f string, args ...interface{}) {
	beego.Error(fmt.Sprintf(f, args...))
}
func (p *serverLogger) Fatalln(args ...interface{}) {
	beego.Error(fmt.Sprintln(args...))
}

func (p *serverLogger) Panic(args ...interface{}) {
	beego.Error(fmt.Sprint(args...))
}
func (p *serverLogger) Panicf(f string, args ...interface{}) {
	beego.Error(fmt.Sprintf(f, args...))
}
func (p *serverLogger) Panicln(args ...interface{}) {
	beego.Error(fmt.Sprintln(args...))
}

var (
	_server     *machinery.Server
	_serverOnce sync.Once
)

// RegisterBackgroundTask register background task
func RegisterBackgroundTask(args ...interface{}) {
	for _, it := range args {
		Server().RegisterTask(GetFunctionName(it), it)
	}
}

// Server background server
func Server() *machinery.Server {
	_serverOnce.Do(func() {
		log.Set(&serverLogger{})
		cfg, err := config.NewFromYaml(path.Join("conf", "server.yml"), false)
		if err != nil {
			beego.Error(err)
			return
		}
		_server, err = machinery.NewServer(cfg)
		if err != nil {
			beego.Error(err)
		}
	})

	return _server
}
