package dict

import (
	"net/http"
	"path/filepath"

	"github.com/chonglou/arche/web/mux"
	"github.com/kapmahc/stardict"
)

func (p *Plugin) postSearch(c *mux.Context) {
	var items []mux.H
	if err := p.dict(func(dt *stardict.Dictionary) error {
		senses := dt.Translate(c.Query("keywords"))
		for _, seq := range senses {
			for _, pat := range seq.Parts {
				items = append(items, mux.H{"type": string(pat.Type), "data": string(pat.Data)})
			}
		}
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (p *Plugin) index(c *mux.Context) {
	items := make(map[string]uint64)
	if err := p.dict(func(d *stardict.Dictionary) error {
		items[d.GetBookName()] = d.GetWordCount()
		return nil
	}); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (p *Plugin) dict(fn func(*stardict.Dictionary) error) error {
	dict, err := stardict.Open(filepath.Join("tmp", "dict"))
	if err != nil {
		return err
	}
	for _, it := range dict {
		if err = fn(it); err != nil {
			return err
		}
	}
	return nil
}
