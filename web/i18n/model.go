package i18n

import (
	"time"
)

// Model locale database model
type Model struct {
	tableName struct{}  `sql:"locales"`
	ID        uint      `json:"id"`
	Lang      string    `json:"lang"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
