package nut

import (
	log "github.com/sirupsen/logrus"
)

// NewHomePage new home page
func NewHomePage() *HomePage {
	return &HomePage{
		handlers: make(map[string]HTMLHandlerFunc),
	}
}

// HomePage home page
type HomePage struct {
	handlers map[string]HTMLHandlerFunc
}

// Register register handler
func (p *HomePage) Register(n string, v HTMLHandlerFunc) {
	if _, ok := p.handlers[n]; ok {
		log.Warnf("handle %s already exist, will override it", n)
	}
	p.handlers[n] = v
}

// Get get handler
func (p *HomePage) Get(name string) HTMLHandlerFunc {
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
