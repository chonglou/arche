package i18n

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"golang.org/x/text/language"
)

// New new
func New(langs ...string) (*I18n, error) {
	var items []language.Tag
	for _, l := range langs {
		t, e := language.Parse(l)
		if e != nil {
			return nil, e
		}
		items = append(items, t)
	}
	return &I18n{
		loaders:   make([]Loader, 0),
		matcher:   language.NewMatcher(items),
		languages: langs[:],
	}, nil
}

// I18n i18n
type I18n struct {
	languages []string
	loaders   []Loader
	matcher   language.Matcher
}

// Languages all available languages
func (p *I18n) Languages() []string {
	return p.languages[:]
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
	for _, it := range p.loaders {
		msg, err := it.Get(lang, code)
		if err == nil && msg != "" {
			return msg, nil
		}
	}
	return "", errors.New(lang + "." + code)
}
