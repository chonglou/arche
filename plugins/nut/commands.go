package nut

import (
	"crypto/x509/pkix"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/migration"
	"github.com/astaxie/beego/orm"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
)

func generateSsl(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "openssl")
		return nil
	}
	root := path.Join("tmp", "etc", "ssl", name)

	key, crt, err := CreateCertificate(
		true,
		pkix.Name{
			Country:      []string{c.String("country")},
			Organization: []string{c.String("organization")},
		},
		c.Int("years"),
	)
	if err != nil {
		return err
	}

	fnk := path.Join(root, "key.pem")
	fnc := path.Join(root, "crt.pem")

	fmt.Printf("generate pem file %s\n", fnk)
	err = WritePemFile(fnk, "RSA PRIVATE KEY", key, 0600)
	fmt.Printf("test: openssl rsa -noout -text -in %s\n", fnk)

	if err == nil {
		fmt.Printf("generate pem file %s\n", fnc)
		err = WritePemFile(fnc, "CERTIFICATE", crt, 0444)
		fmt.Printf("test: openssl x509 -noout -text -in %s\n", fnc)
	}
	if err == nil {
		fmt.Printf(
			"verify: diff <(openssl rsa -noout -modulus -in %s) <(openssl x509 -noout -modulus -in %s)",
			fnk,
			fnc,
		)
	}
	fmt.Println()
	return nil
}

func generateLocale(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "locale")
		return nil
	}
	lng, err := language.Parse(name)
	if err != nil {
		return err
	}
	const root = "locales"
	if err = os.MkdirAll(root, 0700); err != nil {
		return err
	}
	file := path.Join(root, fmt.Sprintf("%s.ini", lng.String()))
	fmt.Printf("generate file %s\n", file)
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	return err
}

func generateMigration(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "migration")
		return nil
	}
	root := filepath.Join(
		migrationsDir(beego.AppConfig.String("databasedriver")),
		fmt.Sprintf("%s_%s", time.Now().Format(migration.DateFormat), name),
	)
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	for _, v := range []string{"up", "down"} {
		fn := filepath.Join(root, fmt.Sprintf("%s.sql", v))
		fmt.Printf("generate file %s\n", fn)
		fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		defer fd.Close()
	}
	return nil
}

func generateNginxConf(c *cli.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	name := beego.BConfig.ServerName

	fn := path.Join("tmp", name+".conf")
	if err = os.MkdirAll(path.Dir(fn), 0700); err != nil {
		return err
	}
	fmt.Printf("generate file %s\n", fn)
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	tpl, err := template.ParseFiles(path.Join("templates", "nginx.conf"))
	if err != nil {
		return err
	}

	return tpl.Execute(fd, struct {
		Port  int
		Root  string
		Name  string
		Theme string
		Ssl   bool
	}{
		Name:  name,
		Port:  beego.BConfig.Listen.HTTPPort,
		Root:  pwd,
		Ssl:   c.Bool("https"),
		Theme: beego.AppConfig.String("theme"),
	})
}

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

	RegisterCommand(cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Usage:   "generate file template",
		Subcommands: []cli.Command{
			{
				Name:    "nginx",
				Aliases: []string{"ng"},
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "https, s",
						Usage: "https?",
					},
				},
				Usage:  "generate nginx.conf",
				Action: generateNginxConf,
			},
			{
				Name:    "openssl",
				Aliases: []string{"ssl"},
				Usage:   "generate ssl certificates",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name, n",
						Usage: "name",
					},
					cli.StringFlag{
						Name:  "country, c",
						Value: "Earth",
						Usage: "country",
					},
					cli.StringFlag{
						Name:  "organization, o",
						Value: "Mother Nature",
						Usage: "organization",
					},
					cli.IntFlag{
						Name:  "years, y",
						Value: 1,
						Usage: "years",
					},
				},
				Action: generateSsl,
			},
			{
				Name:    "migration",
				Usage:   "generate migration file",
				Aliases: []string{"m"},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name, n",
						Usage: "name",
					},
				},
				Action: generateMigration,
			},
			{
				Name:    "locale",
				Usage:   "generate locale file",
				Aliases: []string{"l"},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name, n",
						Usage: "locale name",
					},
				},
				Action: generateLocale,
			},
		},
	})
}
