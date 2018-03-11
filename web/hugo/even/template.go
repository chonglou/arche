package even

import "github.com/chonglou/arche/web/hugo"

// Template template
type Template struct {
}

// Demo demo url
func (p *Template) Demo() string {
	return "https://themes.gohugo.io/theme/hugo-theme-even/"
}

func init() {
	hugo.Register("even", &Template{})
}
