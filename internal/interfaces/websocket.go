package interfaces

import (
	"context"
	"otus-project/internal/model"
)

// WebSocketService интерфейс для WebSocket сервиса
type WebSocketService interface {
	// StartHub запускает WebSocket хаб
	StartHub(ctx context.Context) error

	// StopHub останавливает WebSocket хаб
	StopHub(ctx context.Context) error

	// BroadcastPost отправляет сообщение о новом посте всем подписчикам
	BroadcastPost(ctx context.Context, post *model.WebSocketPost) error

	// GetHub возвращает WebSocket хаб
	GetHub() *model.WebSocketHub
}
