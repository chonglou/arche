package mux

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
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

// Use use middleware
func (p *Router) Use(args ...HandlerFunc) {
	p.handlers = append(p.handlers, args...)
}

// Listen start server
// cors.New(cors.Options{
// 	AllowedOrigins: ["www.change-me.com"],
// 	AllowedMethods: []string{
// 		http.MethodGet,
// 		http.MethodPost,
// 		http.MethodPatch,
// 		http.MethodPut,
// 		http.MethodDelete,
// 	},
// 	AllowedHeaders:   []string{"Authorization", "X-Requested-With"},
// 	AllowCredentials: true,
// 	Debug:            true,
// })
// csrf.Protect(
// 	secret,
// 	csrf.Path("/"),
// 	csrf.Secure(secure),
// 	csrf.CookieName("csrf"),
// 	csrf.RequestHeader("Authenticity-Token"),
// 	csrf.FieldName("authenticity_token"),
// )
func (p *Router) Listen(port int, cors *cors.Cors, csrf mux.MiddlewareFunc, grace bool) error {
	log.Infof(
		"application starting on http://localhost:%d",
		port,
	)
	var hnd http.Handler
	hnd = p.node
	if cors != nil {
		hnd = cors.Handler(hnd)
	}
	if csrf != nil {
		hnd = csrf(hnd)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: hnd,
	}

	if !grace {
		// for debug mode
		return srv.ListenAndServe()
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Warn("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Warn("server exiting")
	return nil
}

// Walk walk all routes
func (p *Router) Walk(f func(path string, methods ...string) error) error {
	return p.node.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		return f(tpl, methods...)
	})
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
		render:   p.render,
	}
}

// Get http get
func (p *Router) Get(pat string, args ...HandlerFunc) {
	p.add(http.MethodGet, pat, args...)
}

// Post http post
func (p *Router) Post(pat string, args ...HandlerFunc) {
	p.add(http.MethodPost, pat, args...)
}

// Delete http delete
func (p *Router) Delete(pat string, args ...HandlerFunc) {
	p.add(http.MethodDelete, pat, args...)
}

// Put http put
func (p *Router) Put(pat string, args ...HandlerFunc) {
	p.add(http.MethodPut, pat, args...)
}

// Patch http patch
func (p *Router) Patch(pat string, args ...HandlerFunc) {
	p.add(http.MethodPatch, pat, args...)
}

func (p *Router) add(met, pat string, args ...HandlerFunc) {
	handlers := append(p.handlers, args...)
	p.node.HandleFunc(pat, func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s", r.Proto, r.Method, r.URL)
		begin := time.Now()
		ctx := Context{
			Request:  r,
			Writer:   w,
			handlers: handlers,
			render:   p.render,
			payload:  make(H),
			index:    -1,
		}
		ctx.Next()
		log.Infof("%s", time.Now().Sub(begin))
	}).Methods(met)
}
