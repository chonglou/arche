package queue

import (
	"reflect"
	"runtime"

	log "github.com/sirupsen/logrus"
)

var handlers = make(map[string]HandlerFunc)

// Register register handler
func Register(n string, h HandlerFunc) {
	if _, ok := handlers[n]; ok {
		log.Warnf("task handler for %s already exists, will override it", n)
	}
	handlers[n] = h
}

// Handlers handlers info
func Handlers() map[string]string {
	items := make(map[string]string)
	for n, f := range handlers {
		items[n] = runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	}
	return items
}

// Get get handler by name
func Get(n string) (HandlerFunc, bool) {
	h, ok := handlers[n]
	return h, ok
}
