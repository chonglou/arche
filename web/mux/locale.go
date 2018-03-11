package mux

import (
	"golang.org/x/text/language"
)

// Locale detect locale from http.Request, order by [query, cookie, header]
func (p *Context) Locale() string {
	const key = "locale"

	// 1. Check URL arguments.
	lang := p.Query(key)

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = p.Cookie(key)
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := p.Header("Accept-Language")
		if len(al) > 4 {
			lang = al[:5] // Only compare first 5 letters.
		}
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
	}

	// Save language information in cookies.
	// if write {
	// 	http.SetCookie(p.Writer, &http.Cookie{
	// 		Name:     key,
	// 		Value:    lang,
	// 		Path:     "/",
	// 		MaxAge:   math.MaxUint32,
	// 		HttpOnly: false,
	// 		Secure:   false,
	// 	})
	// }
	return lang
}
