package nut

import (
	"fmt"
	"path"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
)

func localeFile(l string) string {
	return path.Join("locales", l+".ini")
}

func loadLocales() error {
	for _, lang := range beego.AppConfig.Strings("languages") {
		beego.Debug("Loading language: " + lang)
		if err := i18n.SetMessage(lang, localeFile(lang)); err != nil {
			return err
		}
	}
	return nil
}

// Tr translate content to target language.
func Tr(lang, code string, args ...interface{}) string {
	key := fmt.Sprintf("//locales/%s/%s", lang, code)
	cm := Cache()
	if msg := cm.Get(key); msg != nil {
		return string(msg.([]byte))
	}

	if msg, err := GetLocale(orm.NewOrm(), lang, code); err == nil {
		cm.Put(key, msg, time.Hour*24)
		return fmt.Sprintf(msg, args...)
	}

	return i18n.Tr(lang, code, args...)
}

func init() {
	beego.AddFuncMap("t", Tr)
}
