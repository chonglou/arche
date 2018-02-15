package nut

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
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

// Th translate content to target language.
func Th(lang, code string, arg interface{}) (string, error) {
	msg, err := getLocaleMessage(lang, code)
	if err != nil {
		msg = i18n.Tr(lang, code)
	}
	tpl, err := template.New("").Parse(msg)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, arg)
	return buf.String(), err
}

// Te translate content to target language.
func Te(lang, code string, args ...interface{}) error {
	if msg, err := getLocaleMessage(lang, code); err == nil {
		return fmt.Errorf(msg, args...)
	}
	return errors.New(i18n.Tr(lang, code, args...))
}

// Tr translate content to target language.
func Tr(lang, code string, args ...interface{}) string {
	if msg, err := getLocaleMessage(lang, code); err == nil {
		return fmt.Sprintf(msg, args...)
	}
	return i18n.Tr(lang, code, args...)
}

func getLocaleMessage(lang, code string) (string, error) {
	key := fmt.Sprintf("//locales/%s/%s", lang, code)
	cm := Cache()
	if msg := cm.Get(key); msg != nil {
		return string(msg.([]byte)), nil
	}
	msg, err := GetLocale(orm.NewOrm(), lang, code)
	if err == nil {
		cm.Put(key, msg, time.Hour*24)
	}
	return msg, err
}

func init() {
	beego.AddFuncMap("t", Tr)
}
