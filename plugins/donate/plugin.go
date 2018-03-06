package donate

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
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{}
}

// Mount register
func (p *Plugin) Mount() error {
	rt := p.Router.Group("/donate")
	rt.GET("/projects", p.Layout.HTML("donate/projects/index", p.getProjects))
	rt.GET("/projects/:id", p.Layout.HTML("donate/projects/show", p.getProject))

	api := p.Router.Group("/api/donate", p.Layout.MustSignInMiddleware)
	api.GET("/projects", p.Layout.JSON(p.indexProjects))
	api.POST("/projects", p.Layout.JSON(p.createProject))
	api.GET("/projects/:id", p.Layout.JSON(p.showProject))
	api.POST("/projects/:id", p.Layout.JSON(p.updateProject))
	api.DELETE("/projects/:id", p.Layout.JSON(p.destroyProject))
	return nil
}

func init() {
	web.Register(&Plugin{})
}
