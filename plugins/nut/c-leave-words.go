package nut

import "github.com/astaxie/beego/orm"

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
