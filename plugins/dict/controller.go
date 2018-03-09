package dict

import (
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/stardict"
)

func (p *Plugin) search(_ string, c *gin.Context) error {
	var items []interface{}
	if err := p.dict(func(dt *stardict.Dictionary) error {
		senses := dt.Translate(c.Query("keywords"))
		for _, seq := range senses {
			for _, p := range seq.Parts {
				switch p.Type {
				case 'h':
				case 'g':
					items = append(items, template.HTML(p.Data))
				default:
					items = append(items, string(p.Data))
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	c.Set("results", items)
	return nil
}

func (p *Plugin) index(l string, c *gin.Context) (interface{}, error) {
	items := make(map[string]uint64)
	if err := p.dict(func(d *stardict.Dictionary) error {
		items[d.GetBookName()] = d.GetWordCount()
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
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
