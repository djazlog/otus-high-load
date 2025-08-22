package feed

import (
	"context"
	feedModel "otus-project/internal/repository/feed/model"
)

// Repository интерфейс для работы с материализованной лентой
type Repository interface {
	// AddToFeed добавляет пост в материализованную ленту пользователя
	AddToFeed(ctx context.Context, userID, postID, authorID, postText string) error

	// GetFeed получает материализованную ленту пользователя
	GetFeed(ctx context.Context, userID string, offset, limit int) ([]*feedModel.MaterializedFeed, error)

	// RemoveFromFeed удаляет пост из материализованной ленты пользователя
	RemoveFromFeed(ctx context.Context, userID, postID string) error

	// CreateJob создает задание на материализацию ленты
	CreateJob(ctx context.Context, job *feedModel.FeedJob) error

	// UpdateJobStatus обновляет статус задания
	UpdateJobStatus(ctx context.Context, jobID, status string, error *string) error

	// GetPendingJobs получает задания со статусом pending
	GetPendingJobs(ctx context.Context, limit int) ([]*feedModel.FeedJob, error)

	// GetFriendsOfUser получает список друзей пользователя
	GetFriendsOfUser(ctx context.Context, userID string) ([]string, error)
}
