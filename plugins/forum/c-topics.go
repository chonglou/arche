package forum

import (
	"net/http"
	"time"

	"github.com/chonglou/arche/plugins/nut"
	"github.com/chonglou/arche/web/i18n"
	"github.com/gin-gonic/gin"
)

// -----------------------------------------------------------------------------

func (p *Plugin) indexTopics(l string, c *gin.Context) (interface{}, error) {
	var items []Topic
	if err := p.DB.Model(&items).
		Column("id", "title", "body", "type", "user_id").
		Order("updated_at DESC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmTopic struct {
	Title   string `json:"title" binding:"required"`
	Body    string `json:"body" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Tags    []uint `json:"tags"`
	Catalog uint   `json:"catalog"`
}

func (p *fmTopic) toTags() []*Tag {
	var tags []*Tag
	for _, id := range p.Tags {
		tags = append(tags, &Tag{ID: id})
	}
	return tags
}

func (p *Plugin) createTopic(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmTopic
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Topic{
		Title:     fm.Title,
		Body:      fm.Body,
		Type:      fm.Type,
		CatalogID: fm.Catalog,
		UserID:    user.ID,
		Tags:      fm.toTags(),
	}
	// FIXME
	if err := p.DB.Insert(&it); err != nil {
		return nil, err
	}

	return it, nil
}

func (p *Plugin) showTopic(l string, c *gin.Context) (interface{}, error) {
	var it = Topic{}
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).Select(); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) canEditTopic(c *gin.Context) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var it = Topic{}
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).Select(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	lang := c.MustGet(i18n.LOCALE).(string)
	if it.UserID != user.ID && !c.MustGet(nut.IsAdmin).(bool) {
		c.String(http.StatusForbidden, p.I18n.T(lang, "errors.not-allow"))
		c.Abort()
		return
	}
	c.Set("topic", &it)
}

func (p *Plugin) updateTopic(l string, c *gin.Context) (interface{}, error) {
	var fm fmTopic
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := c.MustGet("topic").(*Topic)
	it.Title = fm.Title
	it.Type = fm.Type
	it.Body = fm.Body
	it.UpdatedAt = time.Now()
	// TODO

	return gin.H{}, nil
}

func (p *Plugin) destroyTopic(l string, c *gin.Context) (interface{}, error) {
	it := c.MustGet("topic").(*Topic)

	if _, err := p.DB.Model(it).Delete(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
