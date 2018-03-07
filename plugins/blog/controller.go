package blog

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
	return filepath.Join("tmp", "blog")
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
		items["/blog"+path[len(rt):]] = path[len(rt)+1 : len(path)-len(ext)]
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
}

func (p *Plugin) show(c *gin.Context) {
	name := c.Param("name")
	if len(name) <= 1 {
		name = "README.md"
	} else {
		name = name[1:]
	}
	if filepath.Ext(name) == ext {
		p.Layout.HTML("blog/show", func(_ string, data gin.H, _ *gin.Context) error {
			title, body, err := p.readMD(name)
			if err != nil {
				return err
			}
			data["title"] = title
			data["body"] = body
			return nil
		})(c)
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(p.root(), name))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	typ := http.DetectContentType(buf)
	c.Data(http.StatusOK, typ, buf)
}

func (p *Plugin) readMD(f string) (string, string, error) {
	buf, err := ioutil.ReadFile(filepath.Join(p.root(), f))
	if err != nil {
		return "", "", err
	}
	return f[:len(f)], string(buf), nil
}
