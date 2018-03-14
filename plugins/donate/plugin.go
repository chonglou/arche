package donate

import (
	"fmt"

	"github.com/chonglou/arche/plugins/nut"
	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/cache"
	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/mux"
	"github.com/chonglou/arche/web/queue"
	"github.com/chonglou/arche/web/settings"
	"github.com/chonglou/arche/web/storage"
	"github.com/facebookgo/inject"
	"github.com/go-pg/pg"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/urfave/cli"
)

// Plugin plugin
type Plugin struct {
	I18n     *i18n.I18n         `inject:""`
	Cache    cache.Cache        `inject:""`
	Jwt      *web.Jwt           `inject:""`
	Dao      *nut.Dao           `inject:""`
	Queue    queue.Queue        `inject:""`
	Settings *settings.Settings `inject:""`
	Security *web.Security      `inject:""`
	Storage  storage.Storage    `inject:""`
	DB       *pg.DB             `inject:""`
	Router   *mux.Router        `inject:""`
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

func (p *Plugin) sitemap() ([]stm.URL, error) {
	var items []stm.URL

	var projects []Project
	if err := p.DB.Model(&projects).
		Column("id", "updated_at").
		Select(); err != nil {
		return nil, err
	}
	for _, it := range projects {
		items = append(
			items,
			stm.URL{
				"loc":     fmt.Sprintf("/donate/projects/%d", it.ID),
				"lastmod": it.UpdatedAt,
			},
		)
	}
	items = append(items, stm.URL{"loc": "/donate/projects"})
	return items, nil
}

// Mount register
func (p *Plugin) Mount() error {
	rt := p.Router.Group("/donate", p.Layout.MustSignInMiddleware)
	rt.GET("/projects", p.indexProjects)
	rt.POST("/projects", p.createProject)
	rt.GET("/projects/{id}", p.showProject)
	rt.POST("/projects/{id}", p.canEditProject, p.updateProject)
	rt.DELETE("/projects/{id}", p.canEditProject, p.destroyProject)
	return nil
}

func init() {
	web.Register(&Plugin{})
}
