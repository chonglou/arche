package i18n

import (
	"time"
)

// Model locale database model
type Model struct {
	tableName struct{} `sql:"locales"`
	ID        uint
	Lang      string
	Code      string
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
