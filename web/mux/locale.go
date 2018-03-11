package mux

import (
	"math"
	"net/http"

	"golang.org/x/text/language"
)

// LOCALE key
const LOCALE = "locale"

// LocaleMiddleware detect locale from http.Request
// order by [query, cookie, header]
func LocaleMiddleware(matcher language.Matcher) HandlerFunc {
	return func(c *Context) {
		write := false

		// 1. Check URL arguments.
		lang := c.Query(LOCALE)

		// 2. Get language information from cookies.
		if len(lang) == 0 {
			lang = c.Cookie(LOCALE)
		} else {
			write = true
		}

		// 3. Get language information from 'Accept-Language'.
		if len(lang) == 0 {
			al := c.Header("Accept-Language")
			if len(al) > 4 {
				lang = al[:5] // Only compare first 5 letters.
			}
			write = true
		}

		// 4. Default language is English.
		tag, err := language.Parse(lang)
		if err != nil {
			tag = language.AmericanEnglish
			write = true
		}

		// 5. Check language is available
		tag, _, _ = matcher.Match(tag)
		if lang != tag.String() {
			lang = tag.String()
			write = true
		}

		// Save language information in cookies.
		if write {
			c.SetCookie(&http.Cookie{
				Name:     LOCALE,
				Value:    lang,
				Path:     "/",
				MaxAge:   math.MaxUint32,
				HttpOnly: false,
				Secure:   false,
			})
		}
		c.Set(LOCALE, lang)
	}
}
