package friend

import (
	"context"
	"github.com/pkg/errors"
)

// AddFriend добавляет друга к пользователю.
func (s *serv) AddFriend(ctx context.Context, userId, friendId string) error {
	if userId == "" || friendId == "" {
		return errors.New("id пользователя или друга не может быть пустым")
	}

	return s.friendRepository.AddFriend(ctx, userId, friendId)
}
