package api

import (
	"encoding/json"
	"log"
	"net/http"
	"otus-project/internal/model"
	"otus-project/internal/utils"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене здесь должна быть проверка origin
	},
}

// WebSocketHandler обрабатывает WebSocket соединения
type WebSocketHandler struct {
	hub *model.WebSocketHub
}

// NewWebSocketHandler создает новый WebSocket обработчик
func NewWebSocketHandler(hub *model.WebSocketHub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket обрабатывает WebSocket соединение для канала /post/feed/posted
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Извлекаем токен из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// Проверяем формат Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Валидируем токен и получаем userID
	claims, err := utils.VerifyToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userID := claims.UserId

	// Обновляем соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to WebSocket: %v", err)
		return
	}

	// Создаем WebSocket соединение
	wsConnection := &model.WebSocketConnection{
		ID:     userID, // Используем userID как ID соединения
		UserID: userID,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
	}

	// Регистрируем соединение в хабе
	h.hub.Register <- wsConnection

	// Запускаем горутины для чтения и записи
	go h.writePump(wsConnection, conn)
	go h.readPump(wsConnection, conn)
}

// writePump отправляет сообщения клиенту
func (h *WebSocketHandler) writePump(connection *model.WebSocketConnection, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		h.hub.Unregister <- connection
	}()

	for {
		select {
		case message, ok := <-connection.Send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// readPump читает сообщения от клиента
func (h *WebSocketHandler) readPump(connection *model.WebSocketConnection, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		h.hub.Unregister <- connection
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Обрабатываем входящие сообщения (если нужно)
		var wsMessage model.WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Printf("Error unmarshaling WebSocket message: %v", err)
			continue
		}

		// Здесь можно добавить обработку различных типов сообщений
		log.Printf("Received message from user %s: %s", connection.UserID, wsMessage.Type)
	}
}
