package i18n

import (
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

func (p *I18n) cacheKey(l, c string) string {
	return "locales/" + l + "/" + c
}
