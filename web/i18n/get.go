package i18n

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

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
		return fmt.Sprintf("%s.%s", lang, code)
	}
	return fmt.Sprintf(msg, args...)
}

// E error
func (p *I18n) E(lang, code string, args ...interface{}) error {
	msg, err := p.get(lang, code)
	if err != nil {
		return fmt.Errorf("%s.%s", lang, code)
	}
	return fmt.Errorf(msg, args...)
}

func (p *I18n) get(lang, code string) (string, error) {
	key := p.cacheKey(lang, code)
	var msg string
	if err := p.cache.Get(key, &msg); err == nil {
		return msg, nil
	}
	if err := p.db.Model(&Model{}).Column("message").
		Where("lang = ?", lang).
		Where("code = ?", code).
		Select(&msg); err != nil {
		return "", err
	}
	p.cache.Put(key, msg, time.Hour*24)
	return msg, nil
}
