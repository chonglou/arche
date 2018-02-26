package nut

import (
	"encoding/base64"

	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/cache"
	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/queue"
	"github.com/chonglou/arche/web/settings"
	"github.com/chonglou/arche/web/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

// Plugin plugin
type Plugin struct {
	I18n     *i18n.I18n         `inject:""`
	Cache    cache.Cache        `inject:""`
	Jwt      *web.Jwt           `inject:""`
	Queue    queue.Queue        `inject:""`
	Settings *settings.Settings `inject:""`
	Security *web.Security      `inject:""`
	Storage  storage.Storage    `inject:""`
	Sitemap  *web.Sitemap       `inject:""`
	RSS      *web.RSS           `inject:""`
	Router   *gin.Engine        `inject:""`
	DB       *pg.DB             `inject:""`
	Dao      *Dao               `inject:""`
	Layout   *Layout            `inject:""`
}

func init() {

	viper.SetDefault("aws", map[string]interface{}{
		"access_key_id":     "change-me",
		"secret_access_key": "change-me",
		"s3": map[string]interface{}{
			"region": "change-me",
			"bucket": "change-me",
		},
	})

	viper.SetDefault("redis", map[string]interface{}{
		"host": "localhost",
		"port": 6379,
		"db":   6,
	})

	viper.SetDefault("rabbitmq", map[string]interface{}{
		"user":     "guest",
		"password": "guest",
		"host":     "localhost",
		"port":     5672,
		"virtual":  "arche-dev",
		"queue":    "tasks",
	})

	viper.SetDefault("postgresql", map[string]interface{}{
		"host":     "localhost",
		"port":     5432,
		"user":     "postgres",
		"password": "",
		"name":     "arche_dev",
		"sslmode":  "disable",
	})

	viper.SetDefault("server", map[string]interface{}{
		"port":    8080,
		"origins": []string{"www.change-me.com"},
		"name":    "www.change-me.com",
	})

	secret, _ := web.RandomBytes(32)
	viper.SetDefault("secret", base64.StdEncoding.EncodeToString(secret))

	viper.SetDefault("elasticsearch", []string{"http://localhost:9200"})

	// ----------------

	web.Register(&Plugin{})
}
