package nut

import (
	"github.com/astaxie/beego/orm"
)

type fmInstall struct {
	Name     string `form:"name" valid:"Required"`
	Email    string `form:"email" valid:"Email;MaxSize(255)"`
	Password string `form:"password" valid:"Required;MinSize(6);MaxSize(32)"`
}

// PostInstall install
// @router /install [post]
func (p *API) PostInstall() {
	p.JSON(func() (interface{}, error) {
		o := orm.NewOrm()

		if cnt, err := o.QueryTable(new(User)).Count(); err != nil {
			return nil, err
		} else if cnt > 0 {
			return nil, Te(p.Lang, "errors.forbidden")
		}
		var fm fmInstall
		if err := p.BindJSON(&fm); err != nil {
			return nil, err
		}

		o.Begin()
		ip := p.Ctx.Input.IP()
		user, err := addEmailUser(o, ip, p.Lang, fm.Name, fm.Email, fm.Password)
		if err != nil {
			o.Rollback()
			return nil, err
		}
		if err := confirmUser(o, ip, p.Lang, user.ID); err != nil {
			o.Rollback()
			return nil, err
		}
		for _, n := range []string{RoleRoot, RoleAdmin} {
			if err := Apply(
				o,
				ip, p.Lang, user.ID,
				n, DefaultResourceType, DefaultResourceID,
				20, 0, 0,
			); err != nil {
				o.Rollback()
				return nil, err
			}
		}
		o.Commit()
		return H{}, nil
	})
}
