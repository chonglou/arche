package forum

import (
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/chonglou/arche/plugins/nut"
)

// IndexTopics index
// @router /topics [get]
func (p *HTML) IndexTopics() {
	p.HTML("forum/topics/index.html", func() error {
		var items []Topic
		if _, err := orm.NewOrm().
			QueryTable(new(Topic)).
			OrderBy("-updated_at").
			All(&items); err != nil {
			return err
		}
		p.Data["topics"] = items
		return nil
	})
}

// ShowTopic show
// @router /topics/:id [get]
func (p *HTML) ShowTopic() {
	p.HTML("forum/topics/show.html", func() error {
		var it Topic
		if err := orm.NewOrm().
			QueryTable(new(Topic)).
			Filter("id", p.Ctx.Input.Param(":id")).
			One(&it); err != nil {
			return err
		}
		p.Data["topic"] = it
		return nil
	})
}

// -----------------------------------------------------------------------------

// IndexTopics index
// @router /topics [get]
func (p *API) IndexTopics() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		var items []Topic
		o := orm.NewOrm()
		q := o.QueryTable(new(Topic))
		if !nut.Is(o, user.ID, nut.RoleAdmin) {
			q = q.Filter("user_id", user.ID)
		}
		if _, err := q.OrderBy("-updated_at").
			All(&items, "id", "title", "updated_at"); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmTopic struct {
	Title string `json:"title" valid:"Required"`
	Body  string `json:"body" valid:"Required"`
	Type  string `json:"type" valid:"Required"`
}

// CreateTopic create
// @router /topics [post]
func (p *API) CreateTopic() {
	p.JSON(func() (interface{}, error) {
		user := p.MustSignIn()
		var fm fmTopic
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		o := orm.NewOrm()
		cid, err := strconv.Atoi(p.Input().Get("catalog"))
		if err != nil {
			return nil, err
		}
		if _, err := o.Insert(&Topic{
			Title:   fm.Title,
			Body:    fm.Body,
			Type:    fm.Type,
			User:    user,
			Catalog: &Catalog{ID: uint(cid)},
		}); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// ShowTopic show
// @router /topics/:id [get]
func (p *API) ShowTopic() {
	p.JSON(func() (interface{}, error) {
		it, err := p.canEditTopic(orm.NewOrm())
		return it, err
	})
}

func (p *API) canEditTopic(o orm.Ormer) (*Topic, error) {
	user := p.MustSignIn()
	var it Topic
	if err := o.
		QueryTable(new(Topic)).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&it); err != nil {
		return nil, err
	}
	if it.User.ID == user.ID || nut.Is(o, user.ID, nut.RoleAdmin) {
		return &it, nil
	}
	return nil, nut.Te(p.Lang, "errors.forbidden")
}

// UpdateTopic update
// @router /topics/:id [post]
func (p *API) UpdateTopic() {
	p.JSON(func() (interface{}, error) {
		var fm fmTopic
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		o := orm.NewOrm()
		it, err := p.canEditTopic(o)
		if err != nil {
			return nil, err
		}
		it.Title = fm.Title
		it.Body = fm.Body
		it.Type = fm.Type
		it.UpdatedAt = time.Now()
		if _, err := o.Update(&it, "title", "body", "type", "updated_at"); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}

// DestroyTopic destroy
// @router /topics/:id [delete]
func (p *API) DestroyTopic() {
	p.JSON(func() (interface{}, error) {
		o := orm.NewOrm()
		it, err := p.canEditTopic(o)
		if err != nil {
			return nil, err
		}
		if _, err := o.Delete(&it); err != nil {
			return nil, err
		}
		return nut.H{}, nil
	})
}
