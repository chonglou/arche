package i18n

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/chonglou/arche/web/cache"
	"github.com/go-pg/pg"
)

// New load from database and filesystem
func New(db *pg.DB, c cache.Cache) *I18n {
	return &I18n{
		db:    db,
		cache: c,
	}
}

// I18n i18n
type I18n struct {
	db    *pg.DB
	cache cache.Cache
}

// Languages all available languages
func (p *I18n) Languages() ([]string, error) {
	var items []string
	if err := p.db.Model(&Model{}).
		ColumnExpr("DISTINCT lang").
		Select(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// All get all items by lang
func (p *I18n) All(l string) (map[string]string, error) {
	var items []Model
	err := p.db.Model(&items).Column("code", "message").
		Where("lang = ?", l).
		Order("code DESC").
		Select()
	if err != nil {
		return nil, err
	}
	rst := make(map[string]string)
	for _, it := range items {
		rst[it.Code] = it.Message
	}
	return rst, nil
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
	key := "locales/" + lang + "/" + code
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
