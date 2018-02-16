package nut

import (
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// IndexAttachments list attachments
// @router /attachments [get]
func (p *API) IndexAttachments() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		o := orm.NewOrm()
		q := o.QueryTable(new(Attachment))
		if !Is(o, user.ID, RoleAdmin) {
			q = q.Filter("user_id", user.ID)
		}
		var items []Attachment
		if _, err := q.OrderBy("-id").All(&items); err != nil {
			return nil, err
		}
		return items, nil
	})
}

// CreateAttachment create
// @router /attachments [post]
func (p *API) CreateAttachment() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()

		_, file, err := p.GetFile("file")
		if err != nil {
			return nil, err
		}
		beego.Debug(file.Filename, file.Header, file.Size)
		fd, err := file.Open()
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}
		fty, href, err := STORAGE().Save(file.Filename, body, file.Size)
		if err != nil {
			return nil, err
		}

		if _, err := orm.NewOrm().Insert(&Attachment{
			Title:        file.Filename,
			Length:       file.Size,
			MediaType:    fty,
			URL:          href,
			ResourceID:   DefaultResourceID,
			ResourceType: DefaultResourceType,
			User:         user,
		}); err != nil {
			return nil, err
		}

		return H{}, nil
	})
}

// DestroyAttachment destroy
// @router /attachments/:id [delete]
func (p *API) DestroyAttachment() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		var it Attachment
		o := orm.NewOrm()
		if err := o.QueryTable(new(Attachment)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it, "id", "user_id"); err != nil {
			return nil, err
		}
		p.canEditAttachment(o, &it, user.ID)
		_, err := o.Delete(&it)
		return H{}, err
	})
}

func (p *API) canEditAttachment(o orm.Ormer, a *Attachment, u uint) {
	if a.User.ID == u || Is(o, u, RoleAdmin) {
		return
	}
	p.CustomAbort(http.StatusForbidden, Tr(p.Lang, "errors.forbidden"))
}
