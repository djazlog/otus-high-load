package user

import (
	"context"
	"otus-project/internal/model"
)

// Search Get получение пользователей по имени и фамилии
func (s *serv) Search(ctx context.Context, filter *model.UserFilter) ([]*model.UserInfo, error) {
	users, err := s.userRepository.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	return users, nil
}
