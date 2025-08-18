package feed

import (
	"context"
	"otus-project/internal/model"
)

// EventHandler обработчик событий для Feed сервиса
type EventHandler struct {
	feedService Service
}

// NewEventHandler создает новый обработчик событий
func NewEventHandler(feedService Service) *EventHandler {
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
