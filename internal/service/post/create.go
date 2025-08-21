package post

import (
	"context"
	"log"
	"otus-project/internal/model"
	"time"
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

	// Публикуем событие создания поста
	if s.eventBus != nil && id != nil && info.Text != nil && info.AuthorUserId != nil {
		event := &model.PostCreatedEvent{
			PostID:       *id,
			AuthorUserID: *info.AuthorUserId,
			PostText:     *info.Text,
			CreatedAt:    time.Now(),
		}

		// Публикуем событие асинхронно, чтобы не блокировать создание поста
		//go func() {
		if err := s.eventBus.PublishEvent(context.Background(), model.EventTypePostCreated, event); err != nil {
			log.Printf("Error publishing post created event: %v", err)
		}
		//}()
	}

	return id, nil
}
