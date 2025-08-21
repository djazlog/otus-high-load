package queue

import (
	"context"
	"otus-project/internal/model"
)

// Client интерфейс для работы с очередью сообщений
type Client interface {
	// PublishFeedEvent публикует событие ленты для конкретного пользователя
	PublishFeedEvent(ctx context.Context, userID string, event *model.FeedEvent) error

	// PublishFeedUpdateTask публикует задачу обновления ленты
	PublishFeedUpdateTask(ctx context.Context, task *model.FeedUpdateTask) error

	// ConsumeFeedMaterializationTasks потребляет задачи материализации ленты
	ConsumeFeedMaterializationTasks(ctx context.Context, handler func(context.Context, *model.FeedUpdateTask) error) error

	// ConsumeFeedEvents потребляет события ленты по routing key feed.event.{user_id}
	ConsumeFeedEvents(ctx context.Context, handler func(context.Context, string, *model.FeedEvent) error) error

	// Close закрывает соединение
	Close() error
}
