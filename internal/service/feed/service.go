package feed

import (
	"context"
	"otus-project/internal/model"
	feedModel "otus-project/internal/repository/feed/model"
)

// Service интерфейс для работы с отложенной материализацией ленты
type Service interface {
	// ScheduleFeedUpdate планирует обновление ленты для друзей автора поста
	ScheduleFeedUpdate(ctx context.Context, postID, authorID, postText string) error

	// ProcessFeedUpdateTask обрабатывает задачу обновления ленты
	ProcessFeedUpdateTask(ctx context.Context, task *model.FeedUpdateTask) error

	// GetMaterializedFeed получает материализованную ленту пользователя
	GetMaterializedFeed(ctx context.Context, userID string, offset, limit int) ([]*feedModel.MaterializedFeed, error)

	// StartWorker запускает воркер для обработки задач материализации
	StartWorker(ctx context.Context) error

	// StopWorker останавливает воркер
	StopWorker(ctx context.Context) error
}
