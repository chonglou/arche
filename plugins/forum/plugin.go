package forum

// inspired by NodeBB

import (
	"github.com/chonglou/arche/plugins/nut"
	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/cache"
	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/queue"
	"github.com/chonglou/arche/web/settings"
	"github.com/chonglou/arche/web/storage"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/urfave/cli"
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
	DB       *pg.DB             `inject:""`
	Router   *gin.Engine        `inject:""`
	Layout   *nut.Layout        `inject:""`
	HomePage *nut.HomePage      `inject:""`
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{}
}

func init() {
	web.Register(&Plugin{})
}
