package nut

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chonglou/arche/web/mux"
)

func (p *Plugin) indexAttachments(c *mux.Context) {
	var items []Attachment
	if err := p.DB.Model(&items).
		Column("id", "title", "url", "media_type", "user_id").
		Order("updated_at DESC").
		Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (p *Plugin) createAttachments(c *mux.Context) {
	user := c.Get(CurrentUser).(*User)
	file, err := c.FormFile("file")
	if err != nil {
		c.Abort(http.StatusBadRequest, err)
		return
	}
	// log.Printf("%s %+v %d", file.Filename, file.Header, file.Size)
	fd, err := file.Open()
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	body, err := ioutil.ReadAll(fd)
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	fty, href, err := p.Storage.Write(file.Filename, body, file.Size)
	if err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	it := Attachment{
		Title:        file.Filename,
		Length:       file.Size,
		MediaType:    fty,
		URL:          href,
		ResourceID:   DefaultResourceID,
		ResourceType: DefaultResourceType,
		UserID:       user.ID,
		UpdatedAt:    time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (p *Plugin) destroyAttachments(c *mux.Context) {
	user := c.Get(CurrentUser).(*User)
	admin := c.Get(IsAdmin).(bool)
	var it Attachment
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).Select(); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	l := c.Get(mux.LOCALE).(string)
	if it.UserID != user.ID && !admin {
		c.Abort(http.StatusForbidden, p.I18n.E(l, "nut.errors.not-allow"))
		return
	}
	if err := p.DB.Delete(&it); err != nil {
		c.Abort(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mux.H{})
}
