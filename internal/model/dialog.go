package model

import "time"

// DialogMessage представляет сообщение в диалоге
type DialogMessage struct {
	// From Идентификатор пользователя отправителя
	From string
	// To Идентификатор пользователя получателя
	To string
	// Text Текст сообщения
	Text string
	// CreatedAt Время создания сообщения
	CreatedAt time.Time
}
