package mux

import (
	"encoding/json"
	"math"
	"mime/multipart"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-playground/form"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
	validator "gopkg.in/go-playground/validator.v9"
)

// HandlerFunc http handler func
type HandlerFunc func(*Context)

// Context http context
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	validate *validator.Validate
	decoder  *form.Decoder
	render   *render.Render
	matcher  language.Matcher
	handlers []HandlerFunc
	params   map[string]string
	payload  H
	index    uint8
}

// Get value by key
func (p *Context) Get(k string) interface{} {
	return p.payload[k]
}

// Set value by key
func (p *Context) Set(k string, v interface{}) {
	p.payload[k] = v
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (p *Context) Next() {
	for s := uint8(len(p.handlers)); p.index < s; {
		f := p.handlers[p.index]
		p.index++
		log.Debugf("call %s", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		f(p)
	}
}

// Home home url
func (p *Context) Home() string {
	scheme := "http"
	if p.Request.TLS != nil {
		scheme += "s"
	}
	return scheme + "://" + p.Request.Host
}

// ClientIP get remote client ip
func (p *Context) ClientIP() string {
	ip := p.Header("X-Forwarded-For")
	if idx := strings.IndexByte(ip, ','); idx >= 0 {
		ip = ip[0:idx]
	}
	ip = strings.TrimSpace(ip)
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(p.Header("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(p.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// FormFile file upload
func (p *Context) FormFile(n string) (*multipart.FileHeader, error) {
	_, fh, er := p.Request.FormFile(n)
	return fh, er
}

// BindForm bind form
func (p *Context) BindForm(fm interface{}) error {
	if err := p.decoder.Decode(
		fm,
		p.Request.URL.Query(),
	); err != nil {
		return err
	}
	return p.validate.Struct(fm)
}

// BindJSON bind json
func (p *Context) BindJSON(fm interface{}) error {
	if err := json.NewDecoder(p.Request.Body).
		Decode(fm); err != nil {
		return err
	}
	return p.validate.Struct(fm)
}

// Header get header value
func (p *Context) Header(n string) string {
	return p.Request.Header.Get(n)
}

// Param get param from url pattern
func (p *Context) Param(n string) string {
	return p.params[n]
}

// Cookie get cookie value
func (p *Context) Cookie(n string) string {
	if ck, er := p.Request.Cookie(n); er == nil {
		return ck.Value
	}
	return ""
}

// SetCookie set cookie
func (p *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(p.Writer, ck)
}

// Query get query param
func (p *Context) Query(n string) string {
	return p.Request.URL.Query().Get(n)
}

// Abort abort
func (p *Context) Abort(s int, e error) {
	p.index = math.MaxUint8
	log.Error(e.Error())
	p.render.Text(p.Writer, s, e.Error())
}

// JSON render json
func (p *Context) JSON(s int, v interface{}) {
	p.render.JSON(p.Writer, s, v)
}

// XML render xml
func (p *Context) XML(s int, v interface{}) {
	p.render.XML(p.Writer, s, v)
}

// HTML render html
func (p *Context) HTML(s int, t string, v interface{}, o render.HTMLOptions) {
	p.render.HTML(p.Writer, s, t, v, o)
}

// Redirect redirect
func (p *Context) Redirect(s int, u string) {
	http.Redirect(p.Writer, p.Request, u, s)
}
