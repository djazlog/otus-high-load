package model

// WebSocketPost представляет сообщение о посте для WebSocket
type WebSocketPost struct {
	PostID       string `json:"postId"`
	PostText     string `json:"postText"`
	AuthorUserID string `json:"author_user_id"`
}

// WebSocketMessage представляет общую структуру WebSocket сообщения
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// WebSocketConnection представляет WebSocket соединение
type WebSocketConnection struct {
	ID     string
	UserID string
	Send   chan []byte
	Hub    *WebSocketHub
}

// WebSocketHub управляет всеми WebSocket соединениями
type WebSocketHub struct {
	Connections map[string]*WebSocketConnection
	Register    chan *WebSocketConnection
	Unregister  chan *WebSocketConnection
	Broadcast   chan *WebSocketPost
}
