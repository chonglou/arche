package initio

import "github.com/chonglou/arche/web/hugo"

// Template template
type Template struct {
}

// Demo demo url
func (p *Template) Demo() string {
	return "https://themes.gohugo.io/theme/hugo-initio/"
}

func init() {
	hugo.Register("initio", &Template{})
}
