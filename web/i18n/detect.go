package i18n

import (
	"math"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// LOCALE locale key
const LOCALE = "locale"

// Middleware detect locale from http.Request, order by [query, cookie, header]
func (p *I18n) Middleware() (gin.HandlerFunc, error) {
	langs, err := p.Languages()
	if err != nil {
		return nil, err
	}
	var tags []language.Tag
	for _, l := range langs {
		t, e := language.Parse(l)
		if e != nil {
			return nil, e
		}
		tags = append(tags, t)
	}

	matcher := language.NewMatcher(tags)
	if err != nil {
		return nil, err
	}
	return func(c *gin.Context) {
		var write bool

		// 1. Check URL arguments.
		lang := c.Query(LOCALE)

		// 2. Get language information from cookies.
		if len(lang) == 0 {
			var err error
			lang, err = c.Cookie(LOCALE)
			if err != nil {
				write = true
			}
		} else {
			write = true
		}

		// 3. Get language information from 'Accept-Language'.
		if len(lang) == 0 {
			al := c.GetHeader("Accept-Language")
			if len(al) > 4 {
				lang = al[:5] // Only compare first 5 letters.
			}

			write = true
		}

		// 4. Default language is English.
		tag, err := language.Parse(lang)
		if err != nil {
			tag = language.AmericanEnglish
		}

		// 5. Check language is available
		tag, _, _ = matcher.Match(tag)
		if lang != tag.String() {
			lang = tag.String()
			write = true
		}

		// Save language information in cookies.
		if write {
			c.SetCookie(LOCALE, lang, math.MaxUint32, "", "", false, false)
		}

		// set payload
		c.Set(LOCALE, lang)
	}, nil
}
