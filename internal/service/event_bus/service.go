package eventBus

import (
	"context"
)

// EventBus интерфейс для системы событий
type EventBus interface {
	// PublishEvent публикует событие
	PublishEvent(ctx context.Context, eventType string, payload interface{}) error

	// Subscribe подписывается на события определенного типа
	Subscribe(eventType string, handler func(context.Context, interface{}) error) error

	// Start запускает обработку событий
	Start(ctx context.Context) error

	// Stop останавливает обработку событий
	Stop(ctx context.Context) error
}
