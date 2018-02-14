package nut

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/migration"
	"github.com/astaxie/beego/orm"
	"github.com/urfave/cli"
)

func runDb(fn func(name string, created int64) error) error {
	db, err := orm.GetDB()
	if err != nil {
		return err
	}

	if err := loadMigrations(beego.AppConfig.String("databasedriver")); err != nil {
		return err
	}
	if err := checkForSchemaUpdateTable(db, orm.NewOrm().Driver().Type()); err != nil {
		return err
	}
	name, created, err := getLatestMigration(db)
	if err != nil {
		return err
	}
	return fn(name, created)
}

func init() {
	RegisterCommand(cli.Command{
		Name:    "database",
		Aliases: []string{"db"},
		Usage:   "database operations",
		Subcommands: []cli.Command{
			{
				Name:    "migrate",
				Usage:   "migrate the DB to the most recent version available",
				Aliases: []string{"m"},
				Action: func(_ *cli.Context) error {
					return runDb(func(name string, created int64) error {
						return migration.Upgrade(created)
					})
				},
			},
			{
				Name:    "rollback",
				Usage:   "roll back the version by 1",
				Aliases: []string{"r"},
				Action: func(_ *cli.Context) error {
					return runDb(func(name string, created int64) error {
						return migration.Rollback(name)
					})
				},
			},
			{
				Name:    "version",
				Usage:   "dump the migration status for the current DB",
				Aliases: []string{"v"},
				Action: func(_ *cli.Context) error {
					return runDb(func(name string, created int64) error {
						tpl := "%-32s %s\n"
						fmt.Printf(tpl, "NAME", "CREATED AT")
						if created > 0 {
							fmt.Printf(tpl, name, time.Unix(created, 0).Format(time.RFC822))
						}
						return nil
					})
				},
			},
		},
	})
}

func dbMigrate(*cli.Context) error {
	return nil
}
