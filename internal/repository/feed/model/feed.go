package model

import (
	"time"
)

// MaterializedFeed представляет материализованную ленту пользователя
type MaterializedFeed struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	PostID    string    `db:"post_id"`
	AuthorID  string    `db:"author_id"`
	PostText  string    `db:"post_text"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// FeedJob представляет задание на материализацию ленты
type FeedJob struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	PostID    string    `db:"post_id"`
	Status    string    `db:"status"` // pending, processing, completed, failed
	Priority  int       `db:"priority"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Error     *string   `db:"error"`
}
