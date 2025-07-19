package service

import (
	"context"

	"otus-project/internal/model"
)

// UserService интерфейс сервиса пользователей
type UserService interface {
	// Register регистрирует пользователя
	Register(ctx context.Context, info *model.UserInfo) (string, error)
	// Get возвращает информацию о пользователе
	Get(ctx context.Context, id string) (*model.UserInfo, error)
	// Search Get возвращает информацию о пользователе
	Search(ctx context.Context, filter *model.UserFilter) ([]*model.UserInfo, error)
	// Login логинит пользователя
	Login(ctx context.Context, login *model.LoginDto) (*string, error)
}
