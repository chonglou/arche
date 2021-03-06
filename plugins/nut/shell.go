package nut

import (
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/chonglou/arche/web"
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
)

const (
	postgresqlDriverName = "postgres"
)

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:    "users",
			Aliases: []string{"us"},
			Usage:   "users operation",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list all users",
					Action: web.InjectAction(func(c *cli.Context) error {
						var items []User
						if err := p.DB.Model(&items).Column("uid", "email", "name").
							Order("last_sign_in_at DESC").
							Select(); err != nil {
							return err
						}
						f := "%-36s %s<%s>\n"
						fmt.Printf(f, "UID", "NAME", "EMAIL")
						for _, it := range items {
							fmt.Printf(f, it.UID, it.Name, it.Email)
						}
						return nil
					}),
				},
				{
					Name:    "role",
					Aliases: []string{"r"},
					Usage:   "apply/deny role to user",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "role's name",
						},
						cli.StringFlag{
							Name:  "user, u",
							Usage: "user's uid",
						},
						cli.IntFlag{
							Name:  "years, y",
							Value: 10,
							Usage: "years",
						},
						cli.BoolFlag{
							Name:  "deny, d",
							Usage: "deny",
						},
						cli.BoolFlag{
							Name:  "confirm, c",
							Usage: "with confirm user",
						},
					},
					Action: web.InjectAction(func(c *cli.Context) error {
						deny := c.Bool("deny")
						confirm := c.Bool("confirm")
						years := c.Int("years")
						name := c.String("name")
						uid := c.String("user")
						if name == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}
						if uid == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}

						return p.DB.RunInTransaction(func(db *pg.Tx) error {
							var user User
							if err := db.Model(&user).Column("id").
								Where("uid = ?", uid).Select(); err != nil {
								return err
							}

							lang := language.AmericanEnglish.String()
							const ip = "0.0.0.0"
							// confirm ?
							if confirm {
								if er := p.Dao.confirmUser(db, lang, ip, &user); er != nil {
									return er
								}
							}

							role, err := p.Dao.GetRole(db, name, DefaultResourceType, DefaultResourceID)
							if err != nil {
								return err
							}
							if deny {
								if err = p.Dao.Deny(db, user.ID, role.ID); err != nil {
									return err
								}
								if err = p.Dao.AddLog(db, user.ID, ip, lang, "nut.logs.user.deny", role); err != nil {
									return err
								}
							} else {
								if err = p.Dao.Allow(db, user.ID, role.ID, years, 0, 0); err != nil {
									return err
								}
								if err = p.Dao.AddLog(db, user.ID, ip, lang, "nut.logs.user.apply", role); err != nil {
									return err
								}
							}
							return nil
						})

					}),
				},
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate file template",
			Subcommands: []cli.Command{
				{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "generate config file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "environment, e",
							Value: "development",
							Usage: "environment, like: development, production, stage, test...",
						},
					},
					Action: p.generateConfig,
				},
				{
					Name:    "nginx",
					Aliases: []string{"ng"},
					Usage:   "generate nginx.conf",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "https, s",
							Usage: "enable HTTP secure?",
						},
						cli.StringFlag{
							Name:  "name, n",
							Usage: "hostname, like: change-me.com",
						},
					},
					Action: web.ConfigAction(p.generateNginxConf),
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
					Action: p.generateSsl,
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
					Action: p.generateMigration,
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
					Action: p.generateLocale,
				},
			},
		},
		{
			Name:    "cache",
			Aliases: []string{"c"},
			Usage:   "cache operations",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Usage:   "list all cache keys",
					Aliases: []string{"l"},
					Action: web.InjectAction(func(_ *cli.Context) error {
						items, err := p.Cache.Status()
						if err != nil {
							return err
						}
						f := "%-36s %v\n"
						fmt.Printf(f, "KEY", "TTL")
						for k, v := range items {
							fmt.Printf(f, k, v)
						}
						return nil
					}),
				},
				{
					Name:    "clear",
					Usage:   "clear cache items",
					Aliases: []string{"c"},
					Action: web.InjectAction(func(_ *cli.Context) error {
						if err := p.Cache.Clear(); err != nil {
							return err
						}
						fmt.Println("Done.")
						return nil
					}),
				},
			},
		},
		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "database operations",
			Subcommands: []cli.Command{
				{
					Name:    "example",
					Usage:   "scripts example for create database and user",
					Aliases: []string{"e"},
					Action:  web.ConfigAction(p.databaseExample),
				},
				{
					Name:    "migrate",
					Usage:   "migrate the DB to the most recent version available",
					Aliases: []string{"m"},
					Action:  p.databaseRun("up"),
				},
				{
					Name:    "rollback",
					Usage:   "roll back the version by 1",
					Aliases: []string{"r"},
					Action:  p.databaseRun("down"),
				},
				{
					Name:    "version",
					Usage:   "dump the migration status for the current DB",
					Aliases: []string{"v"},
					Action:  p.databaseRun("version"),
				},
				{
					Name:    "connect",
					Usage:   "connect database",
					Aliases: []string{"c"},
					Action:  web.ConfigAction(p.connectDatabase),
				},
				{
					Name:    "create",
					Usage:   "create database",
					Aliases: []string{"n"},
					Action:  web.ConfigAction(p.createDatabase),
				},
				{
					Name:    "drop",
					Usage:   "drop database",
					Aliases: []string{"d"},
					Action:  web.ConfigAction(p.dropDatabase),
				},
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start the app server",
			Action: web.InjectAction(func(_ *cli.Context) error {
				// start task queue listener
				go func() {
					// ----------
					host, err := os.Hostname()
					if err != nil {
						log.Error(err)
					}
					for {
						if err := p.Queue.Launch(host); err != nil {
							log.Error(err)
							time.Sleep(5 * time.Second)
						}
					}
				}()
				// -------
				origins := viper.GetStringSlice("server.origins")
				return p.Router.Listen(
					viper.GetInt("server.port"),
					viper.GetString("env") != web.PRODUCTION,
					origins...,
				)
			}),
		},
		{
			Name:    "routes",
			Aliases: []string{"rt"},
			Usage:   "print out all defined routes",
			Action: web.InjectAction(func(_ *cli.Context) error {
				tpl := "%-16s %s\n"
				fmt.Printf(tpl, "METHODS", "PATH")
				return p.Router.Walk(func(methods []string, pattern string) error {
					fmt.Printf(tpl, strings.Join(methods, ","), pattern)
					return nil
				})
			}),
		},
		{
			Name:  "i18n",
			Usage: "internationalization operations",
			Subcommands: []cli.Command{
				{
					Name:    "sync",
					Usage:   "sync locales from locales to database",
					Aliases: []string{"s"},
					Action: web.InjectAction(func(_ *cli.Context) error {
						return p.I18n.Sync("locales")
					}),
				},
			},
		},
	}
}

// --------------------------------------------

func (p *Plugin) generateNginxConf(c *cli.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	name := c.String("name")
	if name == "" {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	fn := path.Join("tmp", "etc", "nginx", "sites-enabled", name+".conf")
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
		Port int
		Root string
		Name string
		Ssl  bool
	}{
		Name: name,
		Port: viper.GetInt("server.port"),
		Root: pwd,
		Ssl:  c.Bool("https"),
	})
}
func (p *Plugin) generateSsl(c *cli.Context) error {
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
func (p *Plugin) generateLocale(c *cli.Context) error {
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
func (p *Plugin) migrationsDir() string {
	return filepath.Join("db", "migrations")
}
func (p *Plugin) generateMigration(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "migration")
		return nil
	}
	root := p.migrationsDir()
	version := time.Now().Format("20060102150405")
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	for _, n := range []string{"up", "down"} {
		dir := filepath.Join(root, fmt.Sprintf("%s_%s", version, name))
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
		fn := filepath.Join(dir, n+".sql")
		fmt.Printf("generate file %s\n", fn)
		fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		defer fd.Close()
	}
	return nil
}
func (p *Plugin) generateConfig(c *cli.Context) error {
	const fn = "config.toml"
	if _, err := os.Stat(fn); err == nil {
		return fmt.Errorf("file %s already exists", fn)
	}
	fmt.Printf("generate file %s\n", fn)

	viper.Set("env", c.String("environment"))

	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	enc := toml.NewEncoder(fd)
	return enc.Encode(viper.AllSettings())
}

// --------------------------------
func (p *Plugin) databaseExample(_ *cli.Context) error {
	args := viper.GetStringMapString("postgresql")
	fmt.Printf("CREATE USER %s WITH PASSWORD '%s';\n", args["user"], args["password"])
	fmt.Printf("CREATE DATABASE %s WITH ENCODING='UTF8';\n", args["name"])
	fmt.Printf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;\n", args["name"], args["user"])
	return nil
}

func (p *Plugin) databaseSource() string {
	args := viper.GetStringMap("postgresql")
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		args["user"],
		args["password"],
		args["host"],
		args["port"],
		args["name"],
		args["sslmode"],
	)
}

func (p *Plugin) databaseRun(act string) cli.ActionFunc {
	return web.ConfigAction(func(_ *cli.Context) error {
		migrations.SetTableName("schema_migrations")
		var items []migrations.Migration
		root := p.migrationsDir()
		files, err := ioutil.ReadDir(root)
		if err != nil {
			return err
		}
		for _, info := range files {
			if !info.IsDir() {
				return nil
			}
			name := info.Name()
			log.Debugf("find migrations %s", name)
			idx := strings.Index(name, "_")
			if idx == -1 {
				return errors.New("bad migration name")
			}
			ver, er := strconv.ParseInt(name[0:idx], 10, 64)
			if er != nil {
				return er
			}
			up, er := ioutil.ReadFile(filepath.Join(root, name, "up.sql"))
			if er != nil {
				return er
			}
			down, er := ioutil.ReadFile(filepath.Join(root, name, "down.sql"))
			if er != nil {
				return er
			}
			items = append(items, migrations.Migration{
				Version: ver,
				Up: func(db migrations.DB) error {
					_, er := db.Exec(string(up))
					return er
				},
				Down: func(db migrations.DB) error {
					_, er := db.Exec(string(down))
					return er
				},
			})
		}
		db, err := p.openDB()
		if err != nil {
			return err
		}
		return db.RunInTransaction(func(tx *pg.Tx) error {
			ov, nv, err := migrations.RunMigrations(db, items, act)
			fmt.Printf("old version: %d, current version: %d\n", ov, nv)
			return err
		})
	})
}

func (p *Plugin) createDatabase(_ *cli.Context) error {
	args := viper.GetStringMapString("postgresql")
	return web.Shell("psql",
		"-h", args["host"],
		"-p", args["port"],
		"-U", "postgres",
		"-c", fmt.Sprintf(
			"CREATE DATABASE %s WITH ENCODING='UTF8'",
			args["name"],
		),
	)
}
func (p *Plugin) dropDatabase(_ *cli.Context) error {
	args := viper.GetStringMapString("postgresql")
	return web.Shell("psql",
		"-h", args["host"],
		"-p", args["port"],
		"-U", "postgres",
		"-c", fmt.Sprintf("DROP DATABASE %s", args["name"]),
	)
}
func (p *Plugin) connectDatabase(_ *cli.Context) error {
	args := viper.GetStringMapString("postgresql")
	return web.Shell("psql",
		"-h", args["host"],
		"-p", args["port"],
		"-U", args["user"],
		args["name"],
	)
}
