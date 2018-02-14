package nut

import (
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/orm"
)

// AllLocales locale map by lang
func AllLocales(o orm.Ormer, lang string) (map[string]string, error) {
	rst := make(map[string]string)
	// load from file
	cfg, err := goconfig.LoadConfigFile(localeFile(lang))
	if err != nil {
		return nil, err
	}
	for _, s := range cfg.GetSectionList() {
		it, er := cfg.GetSection(s)
		if er != nil {
			return nil, er
		}
		for k, v := range it {
			rst[s+"."+k] = v
		}
	}
	// load from database
	var items []Locale
	if _, err := o.QueryTable(new(Locale)).
		Filter("lang", lang).
		All(&items, "code", "message"); err != nil {
		return nil, err
	}
	for _, it := range items {
		rst[it.Code] = it.Message
	}

	return rst, nil
}

// SetLocale set locale
func SetLocale(o orm.Ormer, lang, code, message string) error {
	var it Locale
	err := o.QueryTable(new(Locale)).
		Filter("lang", lang).
		Filter("code", code).
		One(&it, "id")
	if err == nil {
		it.Message = message
		it.UpdatedAt = time.Now()
		_, err = o.Update(&it, "Message", "UpdatedAt")
	} else if err == orm.ErrNoRows {
		_, err = o.Insert(&Locale{
			Lang:    lang,
			Code:    code,
			Message: message,
		})
	}
	return err
}

// GetLocale get locale message
func GetLocale(o orm.Ormer, lang, code string) (string, error) {
	var it Locale
	if err := o.QueryTable(new(Locale)).
		Filter("lang", lang).
		Filter("code", code).
		One(&it, "message"); err != nil {
		return "", err
	}
	return it.Message, nil
}
