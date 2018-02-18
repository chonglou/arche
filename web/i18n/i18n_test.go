package i18n_test

import (
	"testing"
	"time"

	"github.com/chonglou/arche/web/i18n"
	"github.com/chonglou/arche/web/i18n/db"
	"github.com/chonglou/arche/web/i18n/ini"
	"github.com/go-pg/pg"
	"golang.org/x/text/language"
)

func TestIni(t *testing.T) {
	it, err := ini.New("../../locales")
	if err != nil {
		t.Fatal(err)
	}
	testLoader(t, it)
}

func TestDb(t *testing.T) {
	opt, err := pg.ParseURL("postgres://postgres:@localhost:5432/arche?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	con := pg.Connect(opt)
	con.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}
		t.Logf("%s %s", time.Since(event.StartTime), query)
	})

	testLoader(t, db.New(con))
}

func testLoader(t *testing.T, l i18n.Loader) {
	lang := language.AmericanEnglish.String()
	if items, err := l.All(lang); err == nil {
		for k, v := range items {
			t.Log(k, " = ", v)
		}
	} else {
		t.Fatal(err)
	}

	if items, err := l.Langs(); err == nil {
		t.Log(items)
	} else {
		t.Fatal(err)
	}

	if msg, err := l.Get(lang, "languages.en-US"); err == nil {
		t.Log(msg)
	} else {
		t.Fatal(err)
	}
}
