package interfaces

import (
	"context"
	"otus-project/internal/model"
)

// QueueClient интерфейс для работы с очередью сообщений
type QueueClient interface {
	// PublishFeedEvent публикует событие ленты для конкретного пользователя
	PublishFeedEvent(ctx context.Context, userID string, event *model.FeedEvent) error

	// PublishFeedUpdateTask публикует задачу обновления ленты
	PublishFeedUpdateTask(ctx context.Context, task *model.FeedUpdateTask) error

	// ConsumeFeedMaterializationTasks потребляет задачи материализации ленты
	ConsumeFeedMaterializationTasks(ctx context.Context, handler func(context.Context, *model.FeedUpdateTask) error) error

	// Close закрывает соединение
	Close() error
}
