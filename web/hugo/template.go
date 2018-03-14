package hugo

import log "github.com/sirupsen/logrus"

// Template template
type Template interface {
	Demo() string
}

var drivers = make(map[string]Template)

// Register register template
func Register(n string, t Template) {
	if _, ok := drivers[n]; ok {
		log.Warnf("hugo template %s already exists, will override it", n)
	}
	drivers[n] = t
}

// Walk walk templates
func Walk(f func(string, Template) error) error {
	for n, t := range drivers {
		if e := f(n, t); e != nil {
			return e
		}
	}
	return nil
}
