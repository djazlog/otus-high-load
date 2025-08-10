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
}

type FriendRepository interface {
	// AddFriend добавляет нового друга к пользователю.
	AddFriend(ctx context.Context, userId, friendId string) error

	// Delete удаляет связь между пользователем и другом.
	Delete(ctx context.Context, userId, friendId string) error

	// GetFriends возвращает список друзей пользователя.
	GetFriends(ctx context.Context, userId string) ([]string, error)
}
