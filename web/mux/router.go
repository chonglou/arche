package mux

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// New new router
func New(opt render.Options) *Router {
	template.ParseFiles()
	return &Router{
		node:     mux.NewRouter(),
		handlers: make([]HandlerFunc, 0),
		render:   render.New(opt),
	}
}

// Router http router
type Router struct {
	node     *mux.Router
	handlers []HandlerFunc
	render   *render.Render
}

// Static mount static dir like:  node_modules => /3rd/
func (p *Router) Static(pat, dir string) {
	p.node.PathPrefix(pat).Handler(http.StripPrefix(pat, http.FileServer(http.Dir(dir))))
}

// Group group
func (p *Router) Group(pat string, args ...HandlerFunc) *Router {
	return &Router{
		node:     p.node.PathPrefix(pat).Subrouter().StrictSlash(true),
		handlers: append(p.handlers, args...),
	}
}

func (p *Router) add(met, pat string, args ...HandlerFunc) {
	handlers := append(p.handlers, args...)
	p.node.HandleFunc(pat, func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Request:  r,
			Writer:   w,
			handlers: handlers,
			render:   p.render,
			payload:  make(H),
			index:    0,
		}
		ctx.Next()
	})
}
