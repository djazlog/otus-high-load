package model

import "time"

// PostCreatedEvent событие создания поста
type PostCreatedEvent struct {
	PostID       string    `json:"post_id"`
	AuthorUserID string    `json:"author_user_id"`
	PostText     string    `json:"post_text"`
	CreatedAt    time.Time `json:"created_at"`
}

// EventType типы событий
const (
	EventTypePostCreated = "post.created"
)
