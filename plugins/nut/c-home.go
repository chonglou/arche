package nut

// Home home
// @router / [get]
func (p *HTML) Home() {
	p.TplName = "index.tpl"
}
