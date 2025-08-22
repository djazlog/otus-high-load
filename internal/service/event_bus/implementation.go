package eventBus

import (
	"context"
	"log"
	"sync"
)

type service struct {
	handlers map[string][]func(context.Context, interface{}) error
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewService создает новый Event Bus
func NewService() EventBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &service{
		handlers: make(map[string][]func(context.Context, interface{}) error),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// PublishEvent публикует событие
func (s *service) PublishEvent(ctx context.Context, eventType string, payload interface{}) error {
	s.mu.RLock()
	handlers, exists := s.handlers[eventType]
	s.mu.RUnlock()

	if !exists {
		return nil // Нет подписчиков
	}

	// Выполняем обработчики асинхронно
	for _, handler := range handlers {
		go func(h func(context.Context, interface{}) error) {
			if err := h(ctx, payload); err != nil {
				log.Printf("Error handling event %s: %v", eventType, err)
			}
		}(handler)
	}

	return nil
}

// Subscribe подписывается на события определенного типа
func (s *service) Subscribe(eventType string, handler func(context.Context, interface{}) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers[eventType] = append(s.handlers[eventType], handler)
	return nil
}

// Start запускает обработку событий
func (s *service) Start(ctx context.Context) error {
	// Простая реализация - обработка уже происходит в PublishEvent
	return nil
}

// Stop останавливает обработку событий
func (s *service) Stop(ctx context.Context) error {
	s.cancel()
	return nil
}
