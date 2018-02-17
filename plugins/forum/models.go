package forum

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/chonglou/arche/plugins/nut"
)

// Tag Tag
type Tag struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`

	Topics []*Topic `json:"topics" orm:"reverse(many)"`
}

// TableName table name
func (p *Tag) TableName() string {
	return "forum_tags"
}

// Catalog catalog
type Catalog struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`

	Topics []*Topic `json:"topics" orm:"reverse(many)"`
}

// TableName table name
func (p *Catalog) TableName() string {
	return "forum_catalogs"
}

// Topic topic
type Topic struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`

	Posts   []*Post   `json:"posts" orm:"reverse(many)"`
	User    *nut.User `json:"user" orm:"rel(fk)"`
	Tags    []*Tag    `json:"tags" orm:"rel(m2m)"`
	Catalog *Catalog  `json:"catalog" orm:"rel(fk)"`
}

// TableName table name
func (p *Topic) TableName() string {
	return "forum_topics"
}

// Post post
type Post struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`

	User  *nut.User `json:"user" orm:"rel(fk)"`
	Topic *Topic    `json:"topic" orm:"rel(fk)"`
}

// TableName table name
func (p *Post) TableName() string {
	return "forum_posts"
}

// -----------------------------------------------------------------------------

func init() {
	orm.RegisterModel(new(Tag), new(Catalog), new(Topic), new(Post))
}
