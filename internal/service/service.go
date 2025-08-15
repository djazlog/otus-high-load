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

type PostService interface {
	Create(ctx context.Context, info *model.Post) (*string, error)
	Get(ctx context.Context, offset *float32, limit *float32) (*model.Post, error)
	Feed(ctx context.Context, id string, offset *float32, limit *float32) ([]*model.Post, error)
}

type FriendService interface {
	// AddFriend добавляет друга к пользователю.
	AddFriend(ctx context.Context, userId, friendId string) error

	// DeleteFriend удаляет друга из списка пользователя.
	DeleteFriend(ctx context.Context, userId, friendId string) error
}

type DialogService interface {
	// SendMessage отправляет сообщение в диалог
	SendMessage(ctx context.Context, fromUserId, toUserId string, text string) error
	// GetDialogList возвращает список сообщений диалога между двумя пользователями
	GetDialogList(ctx context.Context, userId1, userId2 string) ([]*model.DialogMessage, error)
}
