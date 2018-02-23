package i18n

import (
	"time"

	"github.com/go-pg/pg"
)

// Set set locale
func (p *I18n) Set(db *pg.Tx, lang, code, message string) error {
	var it Model
	err := db.Model(&it).Column("id").
		Where("lang = ?", lang).
		Where("code = ?", code).Select()
	it.Message = message
	it.UpdatedAt = time.Now()
	if err == nil {
		_, err = db.Model(&it).Column("message", "updated_at").Update()
		return err
	}
	if err == pg.ErrNoRows {
		it.Code = code
		it.Lang = lang
		return db.Insert(&it)
	}
	return err
}
