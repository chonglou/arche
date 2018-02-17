package donate

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// Project project
type Project struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	Methods   string    `json:"methods"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
}

// TableName table name
func (p *Project) TableName() string {
	return "donate_projects"
}

// -----------------------------------------------------------------------------

func init() {
	orm.RegisterModel(new(Project))
}
