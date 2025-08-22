package model

import (
	"time"
)

// FeedEvent представляет событие для обновления ленты
type FeedEvent struct {
	PostID       string    `json:"postId"`
	AuthorUserID string    `json:"authorUserId"`
	PostText     string    `json:"postText"`
	CreatedAt    time.Time `json:"createdAt"`
	EventType    string    `json:"eventType"` // "post_created", "post_updated", "post_deleted"
}

// FeedUpdateTask представляет задачу обновления ленты для конкретного пользователя
type FeedUpdateTask struct {
	UserID    string     `json:"userId"`
	PostID    string     `json:"postId"`
	Event     *FeedEvent `json:"event"`
	Priority  int        `json:"priority"` // Приоритет обработки (1 - высокий, 5 - низкий)
	CreatedAt time.Time  `json:"createdAt"`
}

// QueueConfig конфигурация для очереди сообщений
type QueueConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	VHost    string `json:"vhost"`
}

// FeedMaterializationJob задание на материализацию ленты
type FeedMaterializationJob struct {
	UserID    string    `json:"userId"`
	PostID    string    `json:"postId"`
	JobID     string    `json:"jobId"`
	Status    string    `json:"status"` // "pending", "processing", "completed", "failed"
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Error     string    `json:"error,omitempty"`
}
