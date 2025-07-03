package model

import (
	"time"
)

type User struct {
	// Biography Интересы
	Biography *string

	// Birthdate Дата рождения
	Birthdate *time.Time

	// City Город
	City *string

	// FirstName Имя
	FirstName *string

	// Id Идентификатор пользователя
	Id *string

	// SecondName Фамилия
	SecondName *string

	CreatedAt *time.Time

	UpdatedAt *time.Time
}
