package nut

import (
	"compress/gzip"
	"io"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

var sitemapHandlers []SitemapHandler

// SitemapHandler sitemap handler
type SitemapHandler func() ([]stm.URL, error)

// RegisterSitemapHandler register sitemap handler
func RegisterSitemapHandler(args ...SitemapHandler) {
	sitemapHandlers = append(sitemapHandlers, args...)
}

// SitemapXMLGz write sitemap.xml.gz
func SitemapXMLGz(h string, w io.Writer) error {
	sm := stm.NewSitemap()
	sm.Create()
	sm.SetDefaultHost(h)
	for _, hnd := range sitemapHandlers {
		items, err := hnd()
		if err != nil {
			return err
		}
		for _, it := range items {
			sm.Add(it)
		}
	}
	buf := sm.XMLContent()

	wrt := gzip.NewWriter(w)
	defer wrt.Close()
	wrt.Write(buf)
	return nil
}
