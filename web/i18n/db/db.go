package db

import (
	"time"

	"github.com/chonglou/arche/web/i18n"
	"github.com/go-pg/pg"
)

// Locale locale
type Locale struct {
	tableName struct{} `sql:"locales"`
	ID        uint
	Code      string
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New new ini loader
func New(db *pg.DB) i18n.Loader {
	return &Loader{db: db}
}

// Loader load locales from ini file
type Loader struct {
	db *pg.DB
}

// Get get message by lang and code
func (p *Loader) Get(l, c string) (string, error) {
	var msg string
	err := p.db.Model(&Locale{}).Column("message").
		Where("lang = ?", l).
		Where("code = ?", c).Select(&msg)
	return msg, err
}

// Langs list available languages
func (p *Loader) Langs() ([]string, error) {
	var items []string
	err := p.db.Model(&Locale{}).ColumnExpr("DISTINCT lang").Select(&items)
	return items, err
}

// All get all items by lang
func (p *Loader) All(l string) (map[string]string, error) {
	var items []Locale
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
