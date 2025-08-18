package websocket

import (
	"context"
	"encoding/json"
	"log"
	"otus-project/internal/model"
	"sync"
)

type service struct {
	hub    *model.WebSocketHub
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewService создает новый WebSocket сервис
func NewService() WebSocketService {
	ctx, cancel := context.WithCancel(context.Background())

	hub := &model.WebSocketHub{
		Connections: make(map[string]*model.WebSocketConnection),
		Register:    make(chan *model.WebSocketConnection),
		Unregister:  make(chan *model.WebSocketConnection),
		Broadcast:   make(chan *model.WebSocketPost),
	}

	return &service{
		hub:    hub,
		ctx:    ctx,
		cancel: cancel,
	}
}

// StartHub запускает WebSocket хаб
func (s *service) StartHub(ctx context.Context) error {
	go s.runHub()
	log.Println("WebSocket hub started")
	return nil
}

// StopHub останавливает WebSocket хаб
func (s *service) StopHub(ctx context.Context) error {
	s.cancel()
	log.Println("WebSocket hub stopped")
	return nil
}

// BroadcastPost отправляет сообщение о новом посте всем подписчикам
func (s *service) BroadcastPost(ctx context.Context, post *model.WebSocketPost) error {
	s.hub.Broadcast <- post
	return nil
}

// GetHub возвращает WebSocket хаб
func (s *service) GetHub() *model.WebSocketHub {
	return s.hub
}

// runHub запускает основной цикл хаба
func (s *service) runHub() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case connection := <-s.hub.Register:
			s.mu.Lock()
			s.hub.Connections[connection.ID] = connection
			s.mu.Unlock()
			log.Printf("WebSocket connection registered: %s", connection.ID)

		case connection := <-s.hub.Unregister:
			s.mu.Lock()
			if _, ok := s.hub.Connections[connection.ID]; ok {
				delete(s.hub.Connections, connection.ID)
				close(connection.Send)
			}
			s.mu.Unlock()
			log.Printf("WebSocket connection unregistered: %s", connection.ID)

		case post := <-s.hub.Broadcast:
			// Создаем сообщение согласно AsyncAPI спецификации
			message := model.WebSocketMessage{
				Type:    "post",
				Payload: post,
			}

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling WebSocket message: %v", err)
				continue
			}

			s.mu.RLock()
			for _, connection := range s.hub.Connections {
				select {
				case connection.Send <- messageBytes:
				default:
					// Если канал заблокирован, закрываем соединение
					close(connection.Send)
					delete(s.hub.Connections, connection.ID)
				}
			}
			s.mu.RUnlock()
		}
	}
}
