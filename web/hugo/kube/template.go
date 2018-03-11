package kube

import "github.com/chonglou/arche/web/hugo"

// Template template
type Template struct {
}

// Demo demo url
func (p *Template) Demo() string {
	return "https://themes.gohugo.io/theme/kube/"
}

func init() {
	hugo.Register("kube", &Template{})
}
