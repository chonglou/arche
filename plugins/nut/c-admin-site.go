package nut

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/chonglou/arche/web"
	"github.com/chonglou/arche/web/mux"
	"github.com/chonglou/arche/web/queue"
	"github.com/garyburd/redigo/redis"
	"github.com/go-pg/pg"
	gomail "gopkg.in/gomail.v2"
)

func (p *Plugin) deleteAdminSiteClearCache(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	if err := p.Cache.Clear(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) getAdminSiteHome(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var favicon string
	p.Settings.Get(p.DB, "site.favicon", &favicon)
	var theme string
	p.Settings.Get(p.DB, "site.theme", &theme)
	c.JSON(http.StatusOK, mux.H{
		"favicon": favicon,
		"theme":   theme,
	})
}

type fmSiteHome struct {
	Favicon string `json:"favicon" binding:"required"`
	Theme   string `json:"theme" binding:"required"`
}

func (p *Plugin) postAdminSiteHome(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmSiteHome
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	if err := p.DB.RunInTransaction(func(db *pg.Tx) error {
		for k, v := range map[string]interface{}{
			"site.favicon": fm.Favicon,
			"site.theme":   fm.Theme,
		} {
			if err := p.Settings.Set(db, k, v, false); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInsufficientStorage, err)
		return
	}

	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) getAdminSiteSMTP(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var smtp map[string]interface{}
	if err := p.Settings.Get(p.DB, "site.smtp", &smtp); err == nil {
		delete(smtp, "password")
	} else {
		smtp = map[string]interface{}{
			"host":     "localhost",
			"port":     25,
			"username": "whoami@change-me.com",
		}
	}
	c.JSON(http.StatusOK, smtp)
}

type fmSiteSMTP struct {
	Host                 string `json:"host" binding:"required"`
	Port                 int    `json:"port"`
	Username             string `json:"username" binding:"email"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *Plugin) postAdminSiteSMTP(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmSiteSMTP
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	if err := p.Settings.Set(p.DB, "site.smtp", map[string]interface{}{
		"host":     fm.Host,
		"port":     fm.Port,
		"username": fm.Username,
		"password": fm.Password,
	}, true); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) patchAdminSiteSMTP(c *mux.Context) {
	user, err := p.Layout.IsAdmin(c)
	if err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmSiteSMTP
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", fm.Username)
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", "Hi")
	msg.SetBody("text/html", "This is a test email")

	dia := gomail.NewDialer(
		fm.Host,
		fm.Port,
		fm.Username,
		fm.Password,
	)

	if err := dia.DialAndSend(msg); err != nil {
		c.Abort(http.StatusInsufficientStorage, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

const (
	googleSiteVerification = "google-site-verification"
)

func (p *Plugin) getAdminSiteSeo(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var googleVerifyCode string
	p.Settings.Get(p.DB, googleSiteVerification, &googleVerifyCode)
	c.JSON(http.StatusOK, mux.H{
		"googleVerifyCode": googleVerifyCode,
	})
}

type fmSiteSeo struct {
	GoogleVerifyCode string `json:"googleVerifyCode"`
}

func (p *Plugin) postAdminSiteSeo(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmSiteSeo
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	if err := p.DB.RunInTransaction(func(db *pg.Tx) error {
		for k, v := range map[string]string{
			googleSiteVerification: fm.GoogleVerifyCode,
		} {
			if err := p.Settings.Set(db, k, v, false); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

type fmSiteAuthor struct {
	Email string `json:"email" binding:"email"`
	Name  string `json:"name" binding:"required"`
}

func (p *Plugin) postAdminSiteAuthor(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmSiteAuthor
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	if err := p.Settings.Set(p.DB, "site.author", map[string]string{
		"email": fm.Email,
		"name":  fm.Name,
	}, false); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

type fmSiteInfo struct {
	Title       string `json:"title" binding:"required"`
	Subhead     string `json:"subhead" binding:"required"`
	Keywords    string `json:"keywords" binding:"required"`
	Description string `json:"description" binding:"required"`
	Copyright   string `json:"copyright" binding:"required"`
}

func (p *Plugin) postAdminSiteInfo(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	var fm fmSiteInfo
	if err := c.BindJSON(&fm); err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	l := c.Locale()
	if err := p.DB.RunInTransaction(func(db *pg.Tx) error {
		for k, v := range map[string]string{
			"title":       fm.Title,
			"subhead":     fm.Subhead,
			"keywords":    fm.Keywords,
			"description": fm.Description,
			"copyright":   fm.Copyright,
		} {
			if err := p.I18n.Set(db, l, "site."+k, v); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}

func (p *Plugin) getAdminSiteStatus(c *mux.Context) {
	if _, err := p.Layout.IsAdmin(c); err != nil {
		c.Abort(http.StatusForbidden, err)
		return
	}
	ret := mux.H{
		"queue":  queue.Handlers(),
		"routes": p._routes(),
	}
	var err error

	if ret["os"], err = p._osStatus(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	if ret["network"], err = p._networkStatus(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	if ret["postgresql"], err = p._databaseStatus(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	if ret["redis"], err = p._redisStatus(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ret)
}
func (p *Plugin) _osStatus() (mux.H, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	hu, err := user.Current()
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var ifo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&ifo); err != nil {
		return nil, err
	}
	return mux.H{
		"app author":           fmt.Sprintf("%s <%s>", web.AuthorName, web.AuthorEmail),
		"app licence":          web.Copyright,
		"app version":          fmt.Sprintf("%s(%s)", web.Version, web.BuildTime),
		"app root":             pwd,
		"who-am-i":             fmt.Sprintf("%s@%s", hu.Username, hn),
		"go version":           runtime.Version(),
		"go root":              runtime.GOROOT(),
		"go runtime":           runtime.NumGoroutine(),
		"go last gc":           time.Unix(0, int64(mem.LastGC)).Format(time.ANSIC),
		"os cpu":               runtime.NumCPU(),
		"os ram(free/total)":   fmt.Sprintf("%dM/%dM", ifo.Freeram/1024/1024, ifo.Totalram/1024/1024),
		"os swap(free/total)":  fmt.Sprintf("%dM/%dM", ifo.Freeswap/1024/1024, ifo.Totalswap/1024/1024),
		"go memory(alloc/sys)": fmt.Sprintf("%dM/%dM", mem.Alloc/1024/1024, mem.Sys/1024/1024),
		"os time":              time.Now().Format(time.ANSIC),
		"os arch":              fmt.Sprintf("%s(%s)", runtime.GOOS, runtime.GOARCH),
		"os uptime":            (time.Duration(ifo.Uptime) * time.Second).String(),
		"os loads":             ifo.Loads,
		"os procs":             ifo.Procs,
	}, nil
}
func (p *Plugin) _routes() []mux.H {
	var items []mux.H
	p.Router.Walk(func(methods []string, pattern string) error {
		items = append(items, mux.H{
			"methods": strings.Join(methods, ","),
			"path":    pattern,
		})
		return nil
	})
	return items
}
func (p *Plugin) _networkStatus() (mux.H, error) {
	sts := mux.H{}
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, v := range ifs {
		ips := []string{v.HardwareAddr.String()}
		adrs, err := v.Addrs()
		if err != nil {
			return nil, err
		}
		for _, adr := range adrs {
			ips = append(ips, adr.String())
		}
		sts[v.Name] = ips
	}
	return sts, nil
}

func (p *Plugin) _databaseStatus() (mux.H, error) {
	val := mux.H{
		"pool": p.DB.PoolStats(),
		"url":  p.DB.String(),
	}
	var version string
	if _, err := p.DB.Query(pg.Scan(&version), "select version()"); err != nil {
		return nil, err
	}
	val["version"] = version

	// http://blog.javachen.com/2014/04/07/some-metrics-in-postgresql.html
	var size string
	if _, err := p.DB.Query(pg.Scan(&size), "select pg_size_pretty(pg_database_size('postgres'))"); err != nil {
		return nil, err
	}
	val["size"] = size

	// rows, err := db.Query("select pid,current_timestamp - least(query_start,xact_start) AS runtime,substr(query,1,25) AS current_query from pg_stat_activity where not pid=pg_backend_pid()")
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var pid int
	// 	var ts time.Time
	// 	var qry string
	// 	rows.Scan(&pid, &ts, &qry)
	// 	val[fmt.Sprintf("pid-%d", pid)] = fmt.Sprintf("%s (%v)", ts.Format("15:04:05.999999"), qry)
	// }

	return val, nil
}

func (p *Plugin) _redisStatus() (string, error) {
	c := p.Redis.Get()
	defer c.Close()
	return redis.String(c.Do("INFO"))
}
