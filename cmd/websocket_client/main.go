package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type PostPayload struct {
	PostID       string `json:"postId"`
	PostText     string `json:"postText"`
	AuthorUserID string `json:"author_user_id"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <server_url> <auth_token>")
		fmt.Println("Example: go run main.go ws://localhost:8080 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
		os.Exit(1)
	}

	serverURL := os.Args[1]
	token := os.Args[2]

	// Создаем URL для WebSocket соединения
	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal("Error parsing URL:", err)
	}

	// Добавляем путь для WebSocket канала
	u.Path = "/post/feed/posted"

	// Создаем заголовки с токеном авторизации
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	// Подключаемся к WebSocket серверу
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	log.Printf("Connected to WebSocket server: %s", u.String())

	// Канал для сигналов завершения
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Канал для сообщений
	done := make(chan struct{})

	// Горутина для чтения сообщений
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}

			var wsMessage WebSocketMessage
			if err := json.Unmarshal(message, &wsMessage); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			log.Printf("Received message type: %s", wsMessage.Type)

			if wsMessage.Type == "post" {
				// Преобразуем payload в структуру поста
				payloadBytes, _ := json.Marshal(wsMessage.Payload)
				var post PostPayload
				if err := json.Unmarshal(payloadBytes, &post); err != nil {
					log.Printf("Error unmarshaling post payload: %v", err)
					continue
				}

				log.Printf("New post received:")
				log.Printf("  Post ID: %s", post.PostID)
				log.Printf("  Text: %s", post.PostText)
				log.Printf("  Author: %s", post.AuthorUserID)
			}
		}
	}()

	// Ожидаем сигнала завершения
	select {
	case <-interrupt:
		log.Println("Received interrupt signal, closing connection...")

		// Отправляем сообщение о закрытии соединения
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Printf("Error sending close message: %v", err)
		}

		// Ждем завершения горутины чтения
		<-done
	case <-done:
		log.Println("Connection closed by server")
	}

	log.Println("WebSocket client stopped")
}
