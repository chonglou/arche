package nut

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func (p *Plugin) openRender() (*template.Template, error) {
	var items []string
	if err := filepath.Walk(
		filepath.Join("themes", viper.GetString("server.theme"), "views"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				items = append(items, path)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}
	return template.New("").Funcs(template.FuncMap{
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
		"links": func(lng, loc string, x int) ([]Link, error) {
			var items []Link
			if err := p.DB.Model(&items).
				Column("label", "href").
				Where("lang = ? AND loc = ? AND x = ?", lng, loc, x).
				Order("y ASC").
				Select(); err != nil {
				return nil, err
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
	}).ParseFiles(items...)
}
