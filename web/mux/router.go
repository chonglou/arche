package mux

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// NewRouter new router
func NewRouter(o render.Options) *Router {
	return &Router{
		router:   mux.NewRouter(),
		validate: validator.New(),
		decoder:  form.NewDecoder(),
		render:   render.New(o),
		handlers: make([]HandlerFunc, 0),
	}
}

// Router http router
type Router struct {
	router   *mux.Router
	validate *validator.Validate
	decoder  *form.Decoder
	render   *render.Render
	handlers []HandlerFunc
}

// Use using middleware
func (p *Router) Use(args ...HandlerFunc) {
	p.handlers = append(p.handlers, args...)
}

// GET http GET
func (p *Router) GET(pat string, args ...HandlerFunc) {
	p.add(http.MethodGet, pat, args...)
}

// POST http POST
func (p *Router) POST(pat string, args ...HandlerFunc) {
	p.add(http.MethodPost, pat, args...)
}

// PATCH http PATCH
func (p *Router) PATCH(pat string, args ...HandlerFunc) {
	p.add(http.MethodPatch, pat, args...)
}

// PUT http PUT
func (p *Router) PUT(pat string, args ...HandlerFunc) {
	p.add(http.MethodPut, pat, args...)
}

// DELETE http DELETE
func (p *Router) DELETE(pat string, args ...HandlerFunc) {
	p.add(http.MethodDelete, pat, args...)
}

// Group sub-router
func (p *Router) Group(pat string, args ...HandlerFunc) *Router {
	rt := p.router.PathPrefix(pat).Subrouter().StrictSlash(true)
	return &Router{
		router:   rt,
		validate: p.validate,
		decoder:  p.decoder,
		render:   p.render,
		handlers: append(p.handlers, args...),
	}
}

func (p *Router) add(mat, pat string, args ...HandlerFunc) {
	handlers := append(p.handlers, args...)
	p.router.HandleFunc(pat, func(wrt http.ResponseWriter, req *http.Request) {
		begin := time.Now()
		log.Infof("%s %s %s %s", req.Proto, req.Method, req.RemoteAddr, req.RequestURI)
		ctx := &Context{
			Request:  req,
			Writer:   wrt,
			validate: p.validate,
			decoder:  p.decoder,
			render:   p.render,
			params:   mux.Vars(req),
			payload:  H{},
			handlers: handlers,
			index:    0,
		}
		ctx.Next()
		log.Info(time.Now().Sub(begin))
	}).Methods(mat)
}

// Walk walk routes
func (p *Router) Walk(f func(methods []string, pattern string) error) error {
	return p.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pat, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		mtd, err := route.GetMethods()
		if err != nil || len(mtd) == 0 {
			return nil
		}
		return f(mtd, pat)
	})
}

// Listen listen
func (p *Router) Listen(port int, debug bool, origins ...string) error {
	log.Infof(
		"application starting on http://localhost:%d",
		port,
	)
	var hnd http.Handler = p.router

	hnd = cors.New(cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With"},
		AllowCredentials: true,
		Debug:            debug,
	}).Handler(hnd)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: hnd,
	}
	if debug {
		return srv.ListenAndServe()
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Warn("shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server Shutdown:", err)
	}
	log.Warn("server exiting")
	return nil
}
