package i18n

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/chonglou/arche/web/cache"
)

// New new
func New(c cache.Cache) *I18n {
	return &I18n{loaders: make([]Loader, 0), cache: c}
}

// Loader i18n loader
type Loader interface {
	Get(l, c string) (string, error)
	Langs() ([]string, error)
	All(l string) (map[string]string, error)
}

// I18n i18n
type I18n struct {
	loaders []Loader
	cache   cache.Cache
}

// Register register loader
func (p *I18n) Register(args ...Loader) {
	p.loaders = append(p.loaders, args...)
}

// H html
func (p *I18n) H(lang, code string, obj interface{}) (string, error) {
	msg, err := p.get(lang, code)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("").Parse(msg)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, obj)
	return buf.String(), err
}

// T text
func (p *I18n) T(lang, code string, args ...interface{}) string {
	msg, err := p.get(lang, code)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf(msg, args...)
}

// E error
func (p *I18n) E(lang, code string, args ...interface{}) error {
	msg, err := p.get(lang, code)
	if err != nil {
		return err
	}
	return fmt.Errorf(msg, args...)
}

func (p *I18n) get(lang, code string) (string, error) {
	var msg string
	key := fmt.Sprintf("locales/%s/%s", lang, code)
	if err := p.cache.Get(key, &msg); err == nil {
		return msg, nil
	}
	for _, it := range p.loaders {
		msg, err := it.Get(lang, code)
		if err == nil && msg != "" {
			p.cache.Put(key, msg, time.Hour*24)
			return msg, nil
		}
	}
	return "", errors.New(lang + "." + code)
}
