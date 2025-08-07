package friend

import (
	"context"
	"github.com/pkg/errors"
)

// DeleteFriend удаляет друга из списка пользователя.
func (s *serv) DeleteFriend(ctx context.Context, userId, friendId string) error {
	if userId == "" || friendId == "" {
		return errors.New("id пользователя или друга не может быть пустым")
	}

	return s.friendRepository.Delete(ctx, userId, friendId)
}
