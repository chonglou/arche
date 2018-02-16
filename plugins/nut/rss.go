package nut

import (
	"fmt"
	"io"
	"time"

	"github.com/gorilla/feeds"
)

var rssHandlers []RSSHandler

// RSSHandler rss handler
type RSSHandler func(l string) ([]*feeds.Item, error)

// RegisterRSSHandler register rss handler
func RegisterRSSHandler(args ...RSSHandler) {
	rssHandlers = append(rssHandlers, args...)
}

// RSSAtomXML write to rss atom xml
func RSSAtomXML(host, lang, title, dest string, author *feeds.Author, wrt io.Writer) error {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/?locale=%s", host, lang)},
		Description: dest,
		Author:      author,
		Created:     now,
		Items:       make([]*feeds.Item, 0),
	}
	for _, hnd := range rssHandlers {
		items, err := hnd(lang)
		if err != nil {
			return err
		}
		feed.Items = append(feed.Items, items...)
	}

	return feed.WriteAtom(wrt)
}
