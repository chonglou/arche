package nut

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

// NewHTMLRender open render
func NewHTMLRender(root string, funcs template.FuncMap) (render.HTMLRender, error) {
	rdr := HTMLRender{
		layout: "index.html", //filepath.Join(root, "layout", "index.html"),
		tpl:    make(map[string]*template.Template),
	}
	// read layout
	var layout []string
	if err := filepath.Walk(
		filepath.Join(root, "layout"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				layout = append(layout, path)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}
	// read views
	views := filepath.Join(root, "pages")
	if err := filepath.Walk(views, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files := append(layout, path)
			ext := filepath.Ext(path)
			name := path[len(views)+1 : len(path)-len(ext)]

			tpl, err := template.New(rdr.layout).Funcs(funcs).ParseFiles(files...)
			if err != nil {
				return err
			}
			log.Debugf("find template %s => %s, %v", name, tpl.Name(), files)
			rdr.tpl[name] = tpl
		}
		return nil
	},
	); err != nil {
		return nil, err
	}
	return &rdr, nil
}

// HTMLRender html render
type HTMLRender struct {
	layout string
	tpl    map[string]*template.Template
}

// Instance supply render string
func (p *HTMLRender) Instance(name string, data interface{}) render.Render {
	tpl, ok := p.tpl[name]
	if !ok {
		log.Errorf("cann't find template %s", name)
	}
	return render.HTML{
		Template: tpl,
		Data:     data,
		Name:     p.layout,
	}
}

func (p *Plugin) renderFuncMap() template.FuncMap {
	return template.FuncMap{
		"fmt": fmt.Sprintf,
		"str2htm": func(s string) template.HTML {
			return template.HTML(s)
		},
		"dtf": func(t time.Time) string {
			return t.Format(time.RFC822)
		},
		"t": func(lang, code string, args ...interface{}) string {
			return p.I18n.T(lang, code, args...)
		},
		"assets_css": func(u string) template.HTML {
			return template.HTML(fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%s">`, u))
		},
		"assets_js": func(u string) template.HTML {
			return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, u))
		},
		"links": func(lng, loc string) ([]gin.H, error) {
			var xs []int
			if err := p.DB.Model(new(Link)).
				Where("lang = ?", lng).
				Where("loc = ?", loc).
				ColumnExpr("DISTINCT x").
				Order("x ASC").
				Select(&xs); err != nil {
				return nil, err
			}

			var items []gin.H
			for _, x := range xs {
				var links []Link
				if err := p.DB.Model(&links).
					Column("label", "href").
					Where("lang = ?", lng).
					Where("loc = ?", loc).
					Where("x = ?", x).
					Order("y ASC").
					Select(); err != nil {
					return nil, err
				}
				var children []gin.H
				if len(links) > 1 {
					for _, it := range links[1:] {
						children = append(children, gin.H{"label": it.Label, "href": it.Href})
					}
				}
				items = append(items, gin.H{"label": links[0].Label, "href": links[0].Href, "children": children})
			}

			return items, nil
		},
		"cards": func(lng, loc string) ([]Card, error) {
			var items []Card
			if err := p.DB.Model(&items).
				Column("title", "summary", "type", "action", "logo", "href", "loc").
				Where("lang = ? AND loc = ?", lng, loc).
				Order("sort_order ASC").
				Select(); err != nil {
				return nil, err
			}
			return items, nil
		},
		"odd": func(v int) bool {
			return v%2 != 0
		},
		"even": func(v int) bool {
			return v%2 == 0
		},
	}
}
