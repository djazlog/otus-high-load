package user

import (
	"context"
	"otus-project/internal/model"
)

func (s *serv) Login(ctx context.Context, dto *model.LoginDto) (*string, error) {
	user, err := s.userRepository.Login(ctx, dto)
	if err != nil {
		return nil, err
	}

	return user, nil
}
