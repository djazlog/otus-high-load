package model

import (
	"time"
)

type Friend struct {

	// Id Идентификатор пользователя
	UserId *string

	// FriendId Идентификатор пользователя в друзьях
	FriendId *string

	CreatedAt *time.Time

	UpdatedAt *time.Time
}
