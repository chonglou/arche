package ananke

import "github.com/chonglou/arche/web/hugo"

// Template template
type Template struct {
}

// Demo demo url
func (p *Template) Demo() string {
	return "https://themes.gohugo.io/theme/gohugo-theme-ananke/"
}

func init() {
	hugo.Register("ananke", &Template{})
}
