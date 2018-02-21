package nut

import (
	"github.com/chonglou/arche/web/queue"
	"github.com/garyburd/redigo/redis"
)

// Plugin plugin
type Plugin struct {
	I18n      *web.I18n     `inject:""`
	Cache     *web.Cache    `inject:""`
	Jwt       *web.Jwt      `inject:""`
	Jobber    *queue.Queue  `inject:""`
	Settings  *web.Settings `inject:""`
	Security  *web.Security `inject:""`
	S3        *web.S3       `inject:""`
	Sitemap   *web.Sitemap  `inject:""`
	RSS       *web.RSS      `inject:""`
	Router    *gin.Engine   `inject:""`
	DB        *gorm.DB      `inject:""`
	Redis     *redis.Pool   `inject:""`
	Dao       *Dao          `inject:""`
	Layout    *Layout       `inject:""`
	Languages []string      `inject:"languages"`
}
