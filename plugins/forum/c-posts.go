package forum

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chonglou/arche/plugins/nut"
	"github.com/chonglou/arche/web/i18n"
	"github.com/gin-gonic/gin"
)

// -----------------------------------------------------------------------------
func (p *Plugin) indexPosts(l string, c *gin.Context) (interface{}, error) {
	var items []Post
	if err := p.DB.Model(&items).
		Column("id", "body", "type", "updated_at", "topic_id").
		Order("updated_at DESC").Select(); err != nil {
		return nil, err
	}
	return items, nil
}

type fmPost struct {
	Body string `json:"body" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (p *Plugin) createPost(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmPost
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	tid, err := strconv.Atoi(c.Query("topicId"))
	if err != nil {
		return nil, err
	}
	it := Post{
		Type:      fm.Type,
		Body:      fm.Body,
		TopicID:   uint(tid),
		UserID:    user.ID,
		UpdatedAt: time.Now(),
	}
	if err := p.DB.Insert(&it); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showPost(l string, c *gin.Context) (interface{}, error) {
	var it = Post{}
	if err := p.DB.Model(&it).Where("id = ?", c.Param("id")).Select(); err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) canEditPost(c *gin.Context) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var it = Post{}
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
	c.Set("post", &it)
}

func (p *Plugin) updatePost(l string, c *gin.Context) (interface{}, error) {
	var fm fmPost
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := c.MustGet("post").(*Post)
	it.Body = fm.Body
	it.Type = fm.Type
	it.UpdatedAt = time.Now()

	if _, err := p.DB.Model(&it).Column("body", "type", "updated_at").Update(); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyPost(l string, c *gin.Context) (interface{}, error) {
	it := c.MustGet("post").(*Post)
	if err := p.DB.Delete(it); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
