package model

import (
	"time"
)

type LoginDto struct {
	// Id Идентификатор пользователя
	Id string
	// Password Пароль
	Password string
}

type UserInfo struct {
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
	// Password Пароль
	Password *string
	// CreatedAt Дата создания
	CreatedAt *time.Time
	// UpdatedAt Дата обновления
	UpdatedAt *time.Time
}
