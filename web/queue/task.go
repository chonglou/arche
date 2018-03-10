package queue

import (
	"time"

	"github.com/google/uuid"
)

// HandlerFunc task handler func
type HandlerFunc func(id string, body []byte) error

// NewTask create task
func NewTask(t string, p uint8, b []byte) *Task {
	return &Task{
		ID:       uuid.New().String(),
		Type:     t,
		Priority: p,
		Body:     b,
		Created:  time.Now(),
	}
}

// Task task model
type Task struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Body     []byte    `json:"body"`
	Priority uint8     `json:"priority"`
	Created  time.Time `json:"created"`
}
