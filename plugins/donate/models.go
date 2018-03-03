package donate

import (
	"time"

	"github.com/chonglou/arche/plugins/nut"
)

// Project project
type Project struct {
	tableName struct{} `sql:"donate_projects"`

	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	Methods   string    `json:"methods"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    uint      `json:"userId"`
	User      *nut.User `json:"user"`
}
