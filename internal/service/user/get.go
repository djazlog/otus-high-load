package user

import (
	"context"
	"otus-project/internal/model"
)

// Get получение пользователя по id
func (s *serv) Get(ctx context.Context, id string) (*model.UserInfo, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
