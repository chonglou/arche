package i18n

import (
	"math"
	"net/http"

	"golang.org/x/text/language"
)

// LOCALE locale key
const LOCALE = "locale"

// Detect detect locale from http.Request, order by [query, cookie, header]
func Detect(w http.ResponseWriter, r *http.Request) string {
	var write bool

	// 1. Check URL arguments.
	lang := r.URL.Query().Get(LOCALE)

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		if ck, err := r.Cookie(LOCALE); err == nil {
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
	if len(lang) == 0 {
		lang = language.English.String()
		write = true
	}

	// Save language information in cookies.
	if write {
		http.SetCookie(w, &http.Cookie{
			Name:   LOCALE,
			Value:  lang,
			MaxAge: math.MaxUint32,
			Path:   "/",
		})
	}

	return lang
}
