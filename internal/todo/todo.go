package todo

import (
	"time"

	"github.com/google/uuid"
)

type TODO struct {
	ID  uuid.UUID
	Num int32

	Title   string
	Content string

	Completed bool

	CreatedAt time.Time
	DeletedAt *time.Time
}

func New(title string, content string) *TODO {
	return &TODO{
		ID:  uuid.Nil,
		Num: -1,

		Title:   title,
		Content: content,

		Completed: false,

		CreatedAt: time.Now(),
		DeletedAt: nil,
	}
}
