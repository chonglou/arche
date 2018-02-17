package forum

import (
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/chonglou/arche/plugins/nut"
)

// IndexPosts index
// @router /posts [get]
func (p *HTML) IndexPosts() {
	p.HTML("forum/posts/index.html", func() error {
		var items []Post
		if _, err := orm.NewOrm().
			QueryTable(new(Post)).
			OrderBy("-updated_at").
			All(&items); err != nil {
			return err
		}
		p.Data["posts"] = items
		return nil
	})
}

// -----------------------------------------------------------------------------

// IndexPosts index
// @router /posts [get]
func (p *API) IndexPosts() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		var items []Post
		o := orm.NewOrm()
		q := o.QueryTable(new(Topic))
		if !nut.Is(o, user.ID, nut.RoleAdmin) {
			q = q.Filter("user_id", user.ID)
		}
		if _, err := q.OrderBy("-updated_at").
			All(&items, "id", "body", "type", "updated_at"); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmPost struct {
	Body string `json:"body" valid:"Required"`
	Type string `json:"type" valid:"Required"`
}

// CreatePost create
// @router /posts [post]
func (p *API) CreatePost() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		var fm fmPost
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		tid, err := strconv.Atoi(p.Input().Get("topic"))
		if err != nil {
			return nil, err
		}
		if _, err := orm.NewOrm().Insert(&Post{
			Body:  fm.Body,
			Type:  fm.Type,
			Topic: &Topic{ID: uint(tid)},
			User:  user,
		}); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// ShowPost show
// @router /posts/:id [get]
func (p *API) ShowPost() {
	p.JSON(func() (interface{}, error) {
		it, err := p.canEditPost(orm.NewOrm())
		return it, err
	})
}

func (p *API) canEditPost(o orm.Ormer) (*Post, error) {
	user := p.MustSignIn()
	var it Post
	if err := o.
		QueryTable(new(Post)).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&it); err != nil {
		return nil, err
	}
	if it.User.ID == user.ID || nut.Is(o, user.ID, nut.RoleAdmin) {
		return &it, nil
	}
	return nil, nut.Te(p.Lang, "errors.forbidden")
}

// UpdatePost update
// @router /posts/:id [post]
func (p *API) UpdatePost() {
	p.JSON(func() (interface{}, error) {
		var fm fmPost
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		o := orm.NewOrm()
		it, err := p.canEditPost(o)
		if err != nil {
			return nil, err
		}

		it.Body = fm.Body
		it.Type = fm.Type
		it.UpdatedAt = time.Now()
		if _, err := o.Update(&it, "body", "type", "updated_at"); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// DestroyPost destroy
// @router /posts/:id [delete]
func (p *API) DestroyPost() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		o := orm.NewOrm()
		it, err := p.canEditPost(o)
		if err != nil {
			return nil, err
		}
		if _, err := o.Delete(&it); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}
