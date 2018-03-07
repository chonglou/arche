package nut

import (
	log "github.com/sirupsen/logrus"
)

// HomePageHandlerFunc home-page handler func
type HomePageHandlerFunc func(string) (string, error)

// NewHomePage new home page
func NewHomePage() *HomePage {
	return &HomePage{
		handlers: make(map[string]HomePageHandlerFunc),
	}
}

// HomePage home page
type HomePage struct {
	handlers map[string]HomePageHandlerFunc
}

// Register register handler
func (p *HomePage) Register(n string, v HomePageHandlerFunc) {
	if _, ok := p.handlers[n]; ok {
		log.Warnf("handle %s already exist, will override it", n)
	}
	p.handlers[n] = v
}

// Get get handler
func (p *HomePage) Get(name string) HomePageHandlerFunc {
	return p.handlers[name]
}

// Options options
func (p *HomePage) Options() []string {
	var items []string
	for n := range p.handlers {
		items = append(items, n)
	}
	return items
}
