package nut

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	gomail "gopkg.in/gomail.v2"
	yaml "gopkg.in/yaml.v2"
)

// GetAdminSiteHome admin site home
// @router /admin/site/home [get]
func (p *API) GetAdminSiteHome() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		o := orm.NewOrm()
		var favicon string
		Get(o, "site.favicon", &favicon)
		var home map[string]string
		Get(o, "site.home."+p.Lang, &home)
		return H{
			"favicon": favicon,
			"home":    home,
		}, nil
	})
}

type fmSiteHome struct {
	Favicon string `json:"favicon" valid:"Required"`
	Body    string `json:"body" valid:"Required"`
	Type    string `json:"type" valid:"Required"`
}

// PostAdminSiteHome site home
// @router /admin/site/home [post]
func (p *API) PostAdminSiteHome() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmSiteHome
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		o := orm.NewOrm()
		o.Begin()
		for k, v := range map[string]interface{}{
			"site.favicon": fm.Favicon,
			"site.home." + p.Lang: map[string]string{
				"body": fm.Body,
				"type": fm.Type,
			},
		} {
			if err := Set(o, k, v, false); err != nil {
				o.Rollback()
				return nil, err
			}
		}
		o.Commit()
		return H{}, nil
	})
}

// GetAdminSiteSMTP site smtp
// @router /admin/site/smtp [get]
func (p *API) GetAdminSiteSMTP() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var smtp map[string]interface{}
		if err := Get(orm.NewOrm(), "site.smtp", &smtp); err == nil {
			delete(smtp, "password")
		} else {
			smtp = map[string]interface{}{
				"host":     "localhost",
				"port":     25,
				"username": "whoami@change-me.com",
			}
		}
		return smtp, nil
	})
}

type fmSiteSMTP struct {
	Host                 string `json:"host" valid:"Required"`
	Port                 int    `json:"port"`
	Username             string `json:"username" valid:"Email"`
	Password             string `json:"password" valid:"Required"`
	PasswordConfirmation string `json:"passwordConfirmation" valid:"Required"`
}

// PostAdminSiteSMTP site smtp
// @router /admin/site/smtp [post]
func (p *API) PostAdminSiteSMTP() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmSiteSMTP
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if err := Set(orm.NewOrm(), "site.smtp", map[string]interface{}{
			"host":     fm.Host,
			"port":     fm.Port,
			"username": fm.Username,
			"password": fm.Password,
		}, true); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

// PatchAdminSiteSMTP site smtp
// @router /admin/site/smtp [patch]
func (p *API) PatchAdminSiteSMTP() {
	user := p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmSiteSMTP
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
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
			return nil, err
		}
		return H{}, nil
	})
}

const (
	googleSiteVerification = "google-site-verification"
)

// GetAdminSiteSeo site seo
// @router /admin/site/seo [get]
func (p *API) GetAdminSiteSeo() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var googleVerifyCode string
		o := orm.NewOrm()
		Get(o, googleSiteVerification, &googleVerifyCode)
		return H{
			"googleVerifyCode": googleVerifyCode,
		}, nil
	})
}

type fmSiteSeo struct {
	GoogleVerifyCode string `json:"googleVerifyCode"`
}

// PostAdminSiteSeo site seo
// @router /admin/site/seo [post]
func (p *API) PostAdminSiteSeo() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmSiteSeo
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		o := orm.NewOrm()
		o.Begin()
		for k, v := range map[string]string{
			googleSiteVerification: fm.GoogleVerifyCode,
		} {
			if err := Set(o, k, v, false); err != nil {
				o.Rollback()
				return nil, err
			}
		}
		o.Commit()
		return H{}, nil
	})
}

type fmSiteAuthor struct {
	Email string `json:"email" valid:"Email"`
	Name  string `json:"name" valid:"Required"`
}

// PostAdminSiteAuthor site author
// @router /admin/site/author [post]
func (p *API) PostAdminSiteAuthor() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var fm fmSiteAuthor
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if err := Set(orm.NewOrm(), "site.author", map[string]string{
			"email": fm.Email,
			"name":  fm.Name,
		}, false); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

type fmSiteInfo struct {
	Title       string `json:"title" valid:"Required"`
	Subhead     string `json:"subhead" valid:"Required"`
	Keywords    string `json:"keywords" valid:"Required"`
	Description string `json:"description" valid:"Required"`
	Copyright   string `json:"copyright" valid:"Required"`
}

// PostAdminSiteInfo site info
// @router /admin/site/info [post]
func (p *API) PostAdminSiteInfo() {
	p.MustAdmin()

	p.JSON(func() (interface{}, error) {
		var fm fmSiteInfo
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		o := orm.NewOrm()
		o.Begin()
		for k, v := range map[string]string{
			"title":       fm.Title,
			"subhead":     fm.Subhead,
			"keywords":    fm.Keywords,
			"description": fm.Description,
			"copyright":   fm.Copyright,
		} {
			if err := SetLocale(o, p.Lang, "site."+k, v); err != nil {
				o.Rollback()
				return nil, err
			}
		}
		o.Commit()
		return H{}, nil
	})
}

// GetAdminSiteStatus site status
// @router /admin/site/status [get]
func (p *API) GetAdminSiteStatus() {
	p.MustAdmin()

	p.JSON(func() (interface{}, error) {
		ret := H{}
		var err error
		if ret["jobber"], err = p._jobber(); err != nil {
			return nil, err
		}
		if ret["os"], err = p._osStatus(); err != nil {
			return nil, err
		}
		if ret["network"], err = p._networkStatus(); err != nil {
			return nil, err
		}
		if ret["database"], err = p._databaseStatus(); err != nil {
			return nil, err
		}
		if ret["cache"], err = p._cacheStatus(); err != nil {
			return nil, err
		}

		return ret, nil
	})
}

func (p *API) _jobber() (H, error) {
	s := Server()
	cfg, err := yaml.Marshal(s.GetConfig())
	if err != nil {
		return nil, err
	}
	return H{
		"config": string(cfg),
		"tasks":  s.GetRegisteredTaskNames(),
	}, nil
}
func (p *API) _osStatus() (H, error) {
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
	return H{
		"app author":           fmt.Sprintf("%s <%s>", AuthorName, AuthorEmail),
		"app licence":          Copyright,
		"app version":          fmt.Sprintf("%s(%s)", Version, BuildTime),
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
func (p *API) _networkStatus() (H, error) {
	sts := H{}
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

func (p *API) _databaseStatus() (H, error) {
	val := H{
		"drivers": strings.Join(sql.Drivers(), ", "),
	}
	db, err := orm.GetDB()
	if err != nil {
		return nil, err
	}

	switch beego.AppConfig.String("databasedriver") {
	case "postgres":
		var version string
		if err := db.QueryRow("select version()").Scan(&version); err != nil {
			return nil, err
		}
		val["version"] = version

		// http://blog.javachen.com/2014/04/07/some-metrics-in-postgresql.html
		var size string
		if err := db.QueryRow("select pg_size_pretty(pg_database_size('postgres'))").Scan(&size); err != nil {
			return nil, err
		}
		val["size"] = size

		rows, err := db.Query("select pid,current_timestamp - least(query_start,xact_start) AS runtime,substr(query,1,25) AS current_query from pg_stat_activity where not pid=pg_backend_pid()")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var pid int
			var ts time.Time
			var qry string
			rows.Scan(&pid, &ts, &qry)
			val[fmt.Sprintf("pid-%d", pid)] = fmt.Sprintf("%s (%v)", ts.Format("15:04:05.999999"), qry)
		}
	}
	return val, nil
}

func (p *API) _cacheStatus() (string, error) {
	args := make(map[string]string)
	if err := json.Unmarshal(
		[]byte(beego.AppConfig.String("cachesource")),
		&args,
	); err != nil {
		return "", err
	}

	switch beego.AppConfig.String("cachedriver") {
	case "redis":
		pool := &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, e := redis.Dial("tcp", args["conn"])
				if e != nil {
					return nil, e
				}
				if _, e = c.Do("SELECT", args["dbNum"]); e != nil {
					c.Close()
					return nil, e
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
		c := pool.Get()
		defer c.Close()
		return redis.String(c.Do("INFO"))
	default:
		return "", nil
	}
}
