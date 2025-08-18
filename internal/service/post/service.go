package post

import (
	"context"
	"otus-project/internal/client/db"
	"otus-project/internal/interfaces"
	"otus-project/internal/model"
	"otus-project/internal/repository"
	"otus-project/internal/service"
)

type serv struct {
	postPgRepository repository.PostRepository
	postRRepository  repository.PostRepository
	txManager        db.TxManager
	eventBus         interfaces.EventBus
}

// PostService интерфейс сервиса постов
type PostService interface {
	Create(ctx context.Context, info *model.Post) (*string, error)
	Get(ctx context.Context, offset *float32, limit *float32) (*model.Post, error)
	GetByID(ctx context.Context, id string) (*model.Post, error)
	Update(ctx context.Context, id string, text string) error
	Delete(ctx context.Context, id string) error
	Feed(ctx context.Context, id string, offset *float32, limit *float32) ([]*model.Post, error)
}

func NewService(
	postPgRepository repository.PostRepository,
	postRRepository repository.PostRepository,
	txManager db.TxManager,
	eventBus interfaces.EventBus,
) service.PostService {
	return &serv{
		postPgRepository: postPgRepository,
		postRRepository:  postRRepository,
		txManager:        txManager,
		eventBus:         eventBus,
	}
}

// GetByID получает пост по ID
func (s *serv) GetByID(ctx context.Context, id string) (*model.Post, error) {
	return s.postPgRepository.GetByID(ctx, id)
}

// Update обновляет пост
func (s *serv) Update(ctx context.Context, id string, text string) error {
	return s.postPgRepository.Update(ctx, id, text)
}

// Delete удаляет пост
func (s *serv) Delete(ctx context.Context, id string) error {
	return s.postPgRepository.Delete(ctx, id)
}
