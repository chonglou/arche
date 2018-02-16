package nut

import "github.com/astaxie/beego/orm"

// DestroyLeaveWord delete by id
// @router /leave-words/:id [delete]
func (p *API) DestroyLeaveWord() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		if _, err := orm.NewOrm().QueryTable(new(LeaveWord)).
			Filter("id", p.Ctx.Input.Param(":id")).
			Delete(); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}

// IndexLeaveWords list leave-words
// @router /leave-words [get]
func (p *API) IndexLeaveWords() {
	p.MustAdmin()
	p.JSON(func() (interface{}, error) {
		var items []LeaveWord
		if _, err := orm.NewOrm().QueryTable(new(LeaveWord)).
			OrderBy("-created_at").
			All(&items); err != nil {
			return nil, err
		}
		return items, nil
	})
}

type fmLeaveWord struct {
	Type string `json:"type" valid:"Required"`
	Body string `json:"body" valid:"Required"`
}

// CreateLeaveWord create leave-word
// @router /leave-words [post]
func (p *API) CreateLeaveWord() {
	p.JSON(func() (interface{}, error) {
		var fm fmLeaveWord
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}
		if _, err := orm.NewOrm().Insert(&LeaveWord{
			Type: fm.Type,
			Body: fm.Body,
		}); err != nil {
			return nil, err
		}
		return H{}, nil
	})
}
