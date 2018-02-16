package nut

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"github.com/urfave/cli"
)

// RegisterCommand register command
func RegisterCommand(args ...cli.Command) {
	commands = append(commands, args...)
}

var commands []cli.Command

// Main entry
func Main(args ...string) error {

	app := cli.NewApp()
	app.Name = args[0]
	app.Version = fmt.Sprintf("%s (%s) by %s", Version, BuildTime, runtime.Version())
	app.Authors = []cli.Author{
		cli.Author{
			Name:  AuthorName,
			Email: AuthorEmail,
		},
	}
	if ts, err := time.Parse(time.RFC1123Z, BuildTime); err == nil {
		app.Compiled = ts
	}

	app.Copyright = Copyright
	app.Usage = Usage
	app.EnableBashCompletion = true
	app.Commands = commands
	app.Action = func(_ *cli.Context) error {
		// tasks
		toolbox.StartTask()
		defer toolbox.StopTask()

		// start http
		go beego.Run()

		// start worker
		worker := Server().NewWorker(
			beego.BConfig.ServerName,
			beego.AppConfig.DefaultInt("workers", 2),
		)
		if err := worker.Launch(); err != nil {
			beego.Error("worker abort", err)
		}
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	// orm
	orm.Debug = beego.BConfig.RunMode != beego.PROD
	orm.RegisterDataBase(
		"default",
		beego.AppConfig.String("databasedriver"),
		beego.AppConfig.String("databasesource"),
		30,
	)
	// locales
	if err := loadLocales(); err != nil {
		return err
	}

	return app.Run(args)

}
