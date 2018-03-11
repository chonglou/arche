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
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
	validator "gopkg.in/go-playground/validator.v9"
)

// NewRouter new router
func NewRouter(o render.Options, m language.Matcher) *Router {
	return &Router{
		router:   mux.NewRouter(),
		validate: validator.New(),
		decoder:  form.NewDecoder(),
		render:   render.New(o),
		matcher:  m,
	}
}

// Router http router
type Router struct {
	router   *mux.Router
	validate *validator.Validate
	decoder  *form.Decoder
	render   *render.Render
	matcher  language.Matcher
}

// GET http GET
func (p *Router) GET(pat string, hnd HandlerFunc) {
	p.add(http.MethodGet, pat, hnd)
}

// POST http POST
func (p *Router) POST(pat string, hnd HandlerFunc) {
	p.add(http.MethodPost, pat, hnd)
}

// PATCH http PATCH
func (p *Router) PATCH(pat string, hnd HandlerFunc) {
	p.add(http.MethodPatch, pat, hnd)
}

// PUT http PUT
func (p *Router) PUT(pat string, hnd HandlerFunc) {
	p.add(http.MethodPut, pat, hnd)
}

// DELETE http DELETE
func (p *Router) DELETE(pat string, hnd HandlerFunc) {
	p.add(http.MethodDelete, pat, hnd)
}

// Group sub-router
func (p *Router) Group(pat string) *Router {
	rt := p.router.PathPrefix(pat).Subrouter().StrictSlash(true)
	return &Router{
		router:   rt,
		validate: p.validate,
		decoder:  p.decoder,
		render:   p.render,
		matcher:  p.matcher,
	}
}

func (p *Router) add(mat, pat string, hnd HandlerFunc) {
	p.router.HandleFunc(pat, func(wrt http.ResponseWriter, req *http.Request) {
		begin := time.Now()
		log.Info(req.Proto, req.Method, req.RemoteAddr, req.RequestURI)
		hnd(&Context{
			Request:  req,
			Writer:   wrt,
			validate: p.validate,
			decoder:  p.decoder,
			render:   p.render,
			matcher:  p.matcher,
			params:   mux.Vars(req),
		})
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
		if err != nil {
			return err
		}
		if len(mtd) == 0 {
			return nil
		}
		return f(mtd, pat)
	})
}

// Listen listen
func (p *Router) Listen(port int, grace bool) error {
	log.Infof(
		"application starting on http://localhost:%d",
		port,
	)
	var hnd http.Handler = p.router

	// hnd = cors.New(cors.Options{
	// 	AllowedOrigins: viper.GetStringSlice("server.origins"),
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
	// }).Handler(hnd)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: hnd,
	}
	if !grace {
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
