package mux

import (
	"encoding/json"
	"mime/multipart"
	"net"
	"net/http"
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
	params   map[string]string
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
	return p.Param(n)
}

// Cookie get cookie value
func (p *Context) Cookie(n string) string {
	if ck, er := p.Request.Cookie(n); er == nil {
		return ck.Value
	}
	return ""
}

// Query get query param
func (p *Context) Query(n string) string {
	return p.Request.URL.Query().Get(n)
}

// Abort abort
func (p *Context) Abort(s int, e error) {
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
