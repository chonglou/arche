package wiki

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const (
	ext = ".md"
)

func (p *Plugin) root() string {
	return filepath.Join("tmp", "wiki")
}

func (p *Plugin) home(l string, c *gin.Context) error {
	title, body, err := p.readMD("README-" + l + ".md")
	if err != nil {
		return err
	}
	c.Set("title", title)
	c.Set("body", body)
	return nil
}

func (p *Plugin) index(l string, c *gin.Context) (interface{}, error) {
	rt := p.root()
	items := make(map[string]string)
	if err := filepath.Walk(rt, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ext {
			return nil
		}
		items["/wiki"+path[len(rt):]] = path[len(rt)+1 : len(path)-len(ext)]
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
}

func (p *Plugin) show(ctx *gin.Context) {
	name := ctx.Param("name")
	if len(name) <= 1 {
		name = "README.md"
	} else {
		name = name[1:]
	}
	if filepath.Ext(name) == ext {
		p.Layout.HTML("wiki/show", func(_ string, c *gin.Context) error {
			title, body, err := p.readMD(name)
			if err != nil {
				return err
			}
			c.Set("title", title)
			c.Set("body", body)
			return nil
		})(ctx)
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(p.root(), name))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	typ := http.DetectContentType(buf)
	ctx.Data(http.StatusOK, typ, buf)
}

func (p *Plugin) readMD(f string) (string, string, error) {
	buf, err := ioutil.ReadFile(filepath.Join(p.root(), f))
	if err != nil {
		return "", "", err
	}
	return f[:len(f)], string(buf), nil
}
