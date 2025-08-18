package feed

import (
	"context"
	"otus-project/internal/interfaces"
	"otus-project/internal/model"
)

// EventHandler обработчик событий для Feed сервиса
type EventHandler struct {
	feedService interfaces.FeedService
}

// NewEventHandler создает новый обработчик событий
func NewEventHandler(feedService interfaces.FeedService) *EventHandler {
	return &EventHandler{
		feedService: feedService,
	}
}

// HandlePostCreated обрабатывает событие создания поста
func (h *EventHandler) HandlePostCreated(ctx context.Context, payload interface{}) error {
	event, ok := payload.(*model.PostCreatedEvent)
	if !ok {
		return nil // Игнорируем неправильный тип события
	}

	return h.feedService.ScheduleFeedUpdate(ctx, event.PostID, event.AuthorUserID, event.PostText)
}
