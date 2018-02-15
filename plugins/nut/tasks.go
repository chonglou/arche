package nut

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	gomail "gopkg.in/gomail.v2"
)

// SendEmailTask sene email task
func SendEmailTask(to, subject, body string) error {
	beego.Debug("send email to: ", to)
	if beego.BConfig.RunMode != beego.PROD {
		beego.Debug(subject, body)
		return nil
	}

	o := orm.NewOrm()
	smtp := make(map[string]interface{})
	if err := Get(o, "site.smtp", &smtp); err != nil {
		return err
	}

	sender := smtp["username"].(string)
	msg := gomail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dia := gomail.NewDialer(
		smtp["host"].(string),
		smtp["port"].(int),
		sender,
		smtp["password"].(string),
	)

	return dia.DialAndSend(msg)
}

func init() {
	RegisterBackgroundTask(SendEmailTask)
}
