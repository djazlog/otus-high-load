package model

import "time"

// DialogMessage представляет сообщение в диалоге для репозитория
type DialogMessage struct {
	// ID уникальный идентификатор сообщения
	ID int64
	// FromUserID идентификатор пользователя отправителя
	FromUserID string
	// ToUserID идентификатор пользователя получателя
	ToUserID string
	// Text текст сообщения
	Text string
	// CreatedAt время создания сообщения
	CreatedAt time.Time
}
