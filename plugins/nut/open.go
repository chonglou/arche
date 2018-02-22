package nut

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/chonglou/arche/web"
	r_c "github.com/chonglou/arche/web/cache/redis"
	"github.com/chonglou/arche/web/mux"
	"github.com/chonglou/arche/web/queue"
	"github.com/chonglou/arche/web/queue/amqp"
	"github.com/chonglou/arche/web/storage"
	"github.com/chonglou/arche/web/storage/s3"
	"github.com/facebookgo/inject"
	"github.com/garyburd/redigo/redis"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

func (p *Plugin) openDB() (*pg.DB, error) {
	opt, err := pg.ParseURL(p.databaseSource())
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			log.Error(err)
			return
		}
		log.Debugf("%s %s", time.Since(event.StartTime), query)
	})

	return db, nil
}

func (p *Plugin) openS3() (storage.Storage, error) {
	args := viper.GetStringMap("aws")
	s3c := args["s3"].(map[string]string)
	return s3.New(
		args["access_key_id"].(string),
		args["secret_access_key"].(string),
		s3c["region"],
		s3c["bucket"],
	)
}

func (p *Plugin) openQueue() queue.Queue {
	args := viper.GetStringMap("rabbitmq")
	return amqp.New(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		args["user"],
		args["password"],
		args["host"],
		args["port"],
		args["virtual"],
	), args["queue"].(string))
}

func (p *Plugin) openRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(fmt.Sprintf(
				"redis://%s:%d/%d",
				viper.GetString("redis.host"),
				viper.GetInt("redis.port"),
				viper.GetInt("redis.db"),
			))
		},
		// Dial: func() (redis.Conn, error) {
		// 	c, e := redis.Dial(
		// 		"tcp",
		// 		fmt.Sprintf(
		// 			"%s:%d",
		// 			viper.GetString("redis.host"),
		// 			viper.GetInt("redis.port"),
		// 		),
		// 	)
		// 	if e != nil {
		// 		return nil, e
		// 	}
		// 	if _, e = c.Do("SELECT", viper.GetInt("redis.db")); e != nil {
		// 		c.Close()
		// 		return nil, e
		// 	}
		// 	return c, nil
		// },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (p *Plugin) openRender(theme string) render.Options {
	return render.Options{
		Directory:     filepath.Join("themes", theme, "views"),
		Layout:        "layouts/application/index",
		Extensions:    []string{".html"},
		IsDevelopment: web.MODE() != web.PRODUCTION,
		Funcs:         []template.FuncMap{},
	}
}

// Init init beans
func (p *Plugin) Init(g *inject.Graph) error {
	db, err := p.openDB()
	if err != nil {
		return err
	}
	secret, err := base64.StdEncoding.DecodeString(viper.GetString("secret"))
	if err != nil {
		return err
	}

	security, err := web.NewSecurity(secret)
	if err != nil {
		return err
	}

	redis := p.openRedis()

	s3, err := p.openS3()
	if err != nil {
		return err
	}

	theme := viper.GetString("server.theme")

	return g.Provide(
		&inject.Object{Value: db},
		&inject.Object{Value: redis},
		&inject.Object{Value: security},
		&inject.Object{Value: p.openQueue()},
		&inject.Object{Value: s3},
		&inject.Object{Value: web.NewSitemap()},
		&inject.Object{Value: web.NewRSS()},
		&inject.Object{Value: r_c.New(redis, "cache://")},
		&inject.Object{Value: web.NewJwt(secret, crypto.SigningMethodHS512)},
		&inject.Object{Value: mux.New(p.openRender(theme))},
	)
}
