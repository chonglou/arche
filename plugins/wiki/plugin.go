package wiki

// inspired by ArchWiki

import (
	"os"
	"path/filepath"

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
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/urfave/cli"
)

// Plugin plugin(task manager)
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

func (p *Plugin) sitemap() ([]stm.URL, error) {
	var items []stm.URL

	rt := p.root()
	if err := filepath.Walk(rt, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		items = append(items, stm.URL{
			"loc":     "/wiki" + path[len(rt):],
			"lastmod": info.ModTime(),
		})
		return nil
	}); err != nil {
		return nil, err
	}
	items = append(items, stm.URL{"loc": "/wiki/"})
	return items, nil
}

// Mount register
func (p *Plugin) Mount() error {
	p.Sitemap.Register(p.sitemap)
	p.HomePage.Register("wiki/show", p.home)
	// ------------------
	rt := p.Router.Group("/wiki")
	rt.GET("/*name", p.show)

	api := p.Router.Group("/api/wiki")
	api.GET("/", p.Layout.JSON(p.index))
	return nil
}

func init() {
	web.Register(&Plugin{})
}
