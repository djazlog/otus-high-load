package websocket

import (
	"context"
	"otus-project/internal/interfaces"
	"otus-project/internal/model"
)

// EventHandler обработчик событий для WebSocket сервиса
type EventHandler struct {
	websocketService interfaces.WebSocketService
}

// NewEventHandler создает новый обработчик событий
func NewEventHandler(websocketService interfaces.WebSocketService) *EventHandler {
	return &EventHandler{
		websocketService: websocketService,
	}
}

// HandlePostCreated обрабатывает событие создания поста
func (h *EventHandler) HandlePostCreated(ctx context.Context, payload interface{}) error {
	event, ok := payload.(*model.PostCreatedEvent)
	if !ok {
		return nil // Игнорируем неправильный тип события
	}

	wsPost := &model.WebSocketPost{
		PostID:       event.PostID,
		PostText:     event.PostText,
		AuthorUserID: event.AuthorUserID,
	}

	return h.websocketService.BroadcastPost(ctx, wsPost)
}
