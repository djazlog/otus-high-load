package repository

import (
	"context"
	"otus-project/internal/model"
)

type UserRepository interface {
	Login(ctx context.Context, login *model.LoginDto) (*string, error)
	Register(ctx context.Context, info *model.UserInfo) (string, error)
	Get(ctx context.Context, id string) (*model.UserInfo, error)
	Search(ctx context.Context, filter *model.UserFilter) ([]*model.UserInfo, error)
}

type PostRepository interface {
	Create(ctx context.Context, info *model.Post) (*string, error)
	Get(ctx context.Context, offset *float32, limit *float32) (*model.Post, error)
	Feed(ctx context.Context, id string, offset *float32, limit *float32) ([]*model.Post, error)
	CacheFeed(ctx context.Context, userId string, posts []*model.Post) error
	GetByID(ctx context.Context, id string) (*model.Post, error)
	Update(ctx context.Context, id string, text string) error
	Delete(ctx context.Context, id string) error
}

type FriendRepository interface {
	// AddFriend добавляет нового друга к пользователю.
	AddFriend(ctx context.Context, userId, friendId string) error

	// Delete удаляет связь между пользователем и другом.
	Delete(ctx context.Context, userId, friendId string) error

	// GetFriends возвращает список друзей пользователя.
	GetFriends(ctx context.Context, userId string) ([]string, error)
}

type DialogRepository interface {
	// SendMessage сохраняет сообщение в диалоге
	SendMessage(ctx context.Context, fromUserId, toUserId, text string) error
	// GetDialogList возвращает список сообщений диалога между двумя пользователями
	GetDialogList(ctx context.Context, userId1, userId2 string) ([]*model.DialogMessage, error)
}
