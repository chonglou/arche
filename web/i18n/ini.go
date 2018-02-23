package i18n

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-ini/ini"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

const (
	// EXT file ext
	EXT = ".ini"
)

// Sync sync items from filesystem to database
func (p *I18n) Sync(root string) error {
	return p.db.RunInTransaction(func(db *pg.Tx) error {
		return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			name := info.Name()
			if filepath.Ext(name) != EXT {
				return nil
			}

			// detect lang name
			lang := name[0 : len(name)-len(EXT)]
			log.Infof("find locale %s", lang)
			tag, err := language.Parse(lang)
			if err != nil {
				return err
			}
			lang = tag.String()

			// load ini file
			inf, err := ini.Load(path)
			if err != nil {
				return err
			}

			// insert in to database if not exist
			for _, sc := range inf.Sections() {
				for k, msg := range sc.KeysHash() {
					code := sc.Name() + "." + k
					cnt, err := db.Model(new(Model)).
						Where("lang = ?", lang).
						Where("code = ?", code).Count()
					if err != nil {
						return err
					}
					if cnt == 0 {
						if err = db.Insert(&Model{
							Lang:      lang,
							Code:      code,
							Message:   msg,
							UpdatedAt: time.Now(),
						}); err != nil {
							return err
						}
					}
				}
			}

			return nil
		})
	})
}
