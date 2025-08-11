package model

import (
	"time"
)

type Post struct {
	ID           *string
	Text         *string
	AuthorUserId *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}
