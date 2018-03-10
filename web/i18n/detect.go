package i18n

import (
	"math"
	"net/http"

	"golang.org/x/text/language"
)

// Detect detect locale from http.Request, order by [query, cookie, header]
func (p *I18n) Detect(w http.ResponseWriter, r *http.Request) string {
	const key = "locale"
	var write bool

	// 1. Check URL arguments.
	lang := r.URL.Query().Get(key)

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		if ck, err := r.Cookie(key); err == nil {
			lang = ck.Value
		} else {
			write = true
		}
	} else {
		write = true
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := r.Header.Get("Accept-Language")
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
	tag, _, _ = p.matcher.Match(tag)
	if lang != tag.String() {
		lang = tag.String()
		write = true
	}

	// Save language information in cookies.
	if write {
		http.SetCookie(w, &http.Cookie{
			Name:     key,
			Value:    lang,
			Path:     "/",
			MaxAge:   math.MaxUint32,
			HttpOnly: false,
			Secure:   false,
		})
	}
	return lang
}
