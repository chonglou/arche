package nut

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (p *Plugin) indexAttachments(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	admin := c.MustGet(IsAdmin).(bool)
	var items []Attachment
	db := p.DB.Model(&items).Column("id", "title", "url", "media_type")
	if !admin {
		db = db.Where("user_id = ?", user.ID)
	}
	if err := db.Order("updated_at DESC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

func (p *Plugin) createAttachments(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	file, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}
	log.Printf("%s %+v %d", file.Filename, file.Header, file.Size)
	fd, err := file.Open()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	fty, href, err := p.Storage.Write(file.Filename, body, file.Size)
	if err != nil {
		return nil, err
	}
	if err := p.DB.Insert(&Attachment{
		Title:        file.Filename,
		Length:       file.Size,
		MediaType:    fty,
		URL:          href,
		ResourceID:   DefaultResourceID,
		ResourceType: DefaultResourceType,
		UserID:       user.ID,
		UpdatedAt:    time.Now(),
	}); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyAttachments(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	var it Attachment
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).Select(); err != nil {
		return nil, err
	}
	if it.UserID != user.ID && !c.MustGet(IsAdmin).(bool) {
		return nil, p.I18n.E(l, "nut.errors.not-allow")
	}
	if err := p.DB.Delete(&it); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
