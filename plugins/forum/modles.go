package forum

import (
	"time"

	"github.com/chonglou/arche/plugins/nut"
)

// Tag tag
type Tag struct {
	tableName struct{} `sql:"forum_tags"`

	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Topics []*Topic `pg:",many2many:forum_topics_tags,joinFK:topic_id"`
}

// Catalog catalog
type Catalog struct {
	tableName struct{} `sql:"forum_catalogs"`

	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Posts []*Post
}

// Post post
type Post struct {
	tableName struct{} `sql:"forum_posts"`

	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	UserID  uint `json:"userId"`
	User    *nut.User
	TopicID uint `json:"topicId"`
	Topic   *Topic
}

// Topic topic
type Topic struct {
	tableName struct{} `sql:"forum_topics"`

	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	UserID    uint      `json:"userId"`
	User      *nut.User `json:"user"`
	CatalogID uint      `json:"catalogId"`
	Catalog   *Catalog  `json:"catalog"`
	Posts     []*Post   `json:"posts"`

	Tags []*Tag `pg:",many2many:forum_topics_tags,joinFK:tag_id"`
}
