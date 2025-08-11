package post

import (
	"context"
	"otus-project/internal/model"
)

// Create Создание поста
func (s *serv) Create(ctx context.Context, info *model.Post) (*string, error) {
	var id *string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.postPgRepository.Create(ctx, info)
		if errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return id, nil
}
