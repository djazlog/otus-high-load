package post

import (
	"context"
	"otus-project/internal/model"
)

func (s *serv) Get(ctx context.Context, offset *float32, limit *float32) (*model.Post, error) {
	note, err := s.postPgRepository.Get(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return note, nil
}
