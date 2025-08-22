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

type PostList struct {
	Posts []*Post
}

type MaterializedFeed struct {
	ID        string   
	UserID    string    
	PostID    string    
	AuthorID  string    
	PostText  string    
	CreatedAt time.Time 
}
