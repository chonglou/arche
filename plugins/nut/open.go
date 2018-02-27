package nut

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/chonglou/arche/web"
	r_c "github.com/chonglou/arche/web/cache/redis"
	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/queue"
	"github.com/chonglou/arche/web/queue/amqp"
	"github.com/chonglou/arche/web/settings"
	"github.com/chonglou/arche/web/storage"
	"github.com/chonglou/arche/web/storage/fs"
	"github.com/chonglou/arche/web/storage/s3"
	"github.com/facebookgo/inject"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func (p *Plugin) openStorage() (storage.Storage, error) {
	typ := viper.GetString("storage.provider")
	switch typ {
	case "s3":
		args := viper.GetStringMap("aws")
		s3c := args["s3"].(map[string]interface{})
		return s3.New(
			args["access_key_id"].(string),
			args["secret_access_key"].(string),
			s3c["region"].(string),
			s3c["bucket"].(string),
		)
	case "local":
		return fs.New(
			viper.GetString("storage.root"),
			viper.GetString("storage.endpoint"),
		)

	}
	return nil, fmt.Errorf("bad storage provider %s", typ)

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
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
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
	kvs, err := settings.New(secret)
	if err != nil {
		return err
	}

	redis := p.openRedis()

	st, err := p.openStorage()
	if err != nil {
		return err
	}

	cache := r_c.New(redis, "cache://")

	// i18n
	i18n := i18n.New(db, cache)

	return g.Provide(
		&inject.Object{Value: db},
		&inject.Object{Value: redis},
		&inject.Object{Value: security},
		&inject.Object{Value: kvs},
		&inject.Object{Value: p.openQueue()},
		&inject.Object{Value: st},
		&inject.Object{Value: web.NewSitemap()},
		&inject.Object{Value: web.NewRSS()},
		&inject.Object{Value: cache},
		&inject.Object{Value: web.NewJwt(secret, crypto.SigningMethodHS512)},
		&inject.Object{Value: i18n},
		&inject.Object{Value: gin.Default()},
	)
}
