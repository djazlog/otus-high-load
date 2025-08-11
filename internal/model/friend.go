package model

import (
	"time"
)

type Friend struct {
	// Id Идентификатор пользователя
	UserId *string
	// FriendId Идентификатор пользователя в друзьях
	FriendId *string
	// CreatedAt Дата создания
	CreatedAt *time.Time
	// UpdatedAt Дата обновления
	UpdatedAt *time.Time
}
