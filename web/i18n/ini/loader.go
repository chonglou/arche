package ini

import (
	"os"
	"path/filepath"

	"github.com/chonglou/arche/web/i18n"
	"github.com/go-ini/ini"
)

type locales map[string]string

// New new ini loader
func New(d string) (i18n.Loader, error) {
	items := make(map[string]locales)
	if err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
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

		lang := name[0 : len(name)-len(EXT)]
		inf, err := ini.Load(path)
		if err != nil {
			return err
		}
		rst := make(locales)
		for _, sc := range inf.Sections() {
			for k, v := range sc.KeysHash() {
				rst[sc.Name()+"."+k] = v
			}
		}

		items[lang] = rst
		return nil
	}); err != nil {
		return nil, err
	}
	return &Loader{items: items}, nil
}

// Loader load locales from ini file
type Loader struct {
	items map[string]locales
}

// Langs list available languages
func (p *Loader) Langs() ([]string, error) {
	var items []string
	for k, _ := range p.items {
		items = append(items, k)
	}
	return items, nil
}

// Get get message by lang and code
func (p *Loader) Get(l, c string) (string, error) {
	it, ok := p.items[l]
	if !ok {
		return "", nil
	}
	return it[c], nil
}

// All read all items by lang
func (p *Loader) All(l string) (map[string]string, error) {
	rst := make(map[string]string)
	if it, ok := p.items[l]; ok {
		for k, v := range it {
			rst[k] = v
		}
	}
	return rst, nil
}

const (
	// EXT file ext
	EXT = ".ini"
)
