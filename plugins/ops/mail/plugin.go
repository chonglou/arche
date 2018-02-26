package mail

import (
	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/cache"
	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/queue"
	"github.com/chonglou/arche/web/settings"
	"github.com/chonglou/arche/web/storage"
	"github.com/facebookgo/inject"
	"github.com/go-pg/pg"
	"github.com/urfave/cli"
)

// https://wiki2.dovecot.org/HowTo/DovecotPostgresql
// http://www.postfix.org/PGSQL_README.html
// https://linode.com/docs/email/postfix/email-with-postfix-dovecot-and-mysql/
// https://www.tunnelsup.com/using-salted-sha-hashes-with-dovecot-authentication/
// https://wiki2.dovecot.org/Authentication/PasswordSchemes

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
	return nil
}

func init() {
	web.Register(&Plugin{})
}
