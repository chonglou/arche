package mux

import (
	"math"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

// H hash
type H map[string]interface{}

// Context http context
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	handlers []HandlerFunc
	render   *render.Render
	payload  H
	index    int8
}

// Next run next
func (p *Context) Next() {
	p.index++
	for s := int8(len(p.handlers)); p.index < s; p.index++ {
		p.handlers[p.index](p)
	}
}

// Abort abort
func (p *Context) Abort(s int, e error) {
	p.Text(s, e.Error())
	log.Error(e)
	p.index = math.MaxInt8
}

// JSON render json
func (p *Context) JSON(s int, v interface{}) {
	p.render.JSON(p.Writer, s, v)
}

// XML render xml
func (p *Context) XML(s int, v interface{}) {
	p.render.XML(p.Writer, s, v)
}

// Text render text
func (p *Context) Text(s int, v string) {
	p.render.Text(p.Writer, s, v)
}

// HTML render html
func (p *Context) HTML(s int, t string, v interface{}, args ...render.HTMLOptions) {
	p.render.HTML(p.Writer, s, t, v, args...)
}
