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

	// Отправляем уведомление через WebSocket
	if s.websocketService != nil && id != nil && info.Text != nil && info.AuthorUserId != nil {
		wsPost := &model.WebSocketPost{
			PostID:       *id,
			PostText:     *info.Text,
			AuthorUserID: *info.AuthorUserId,
		}

		// Отправляем асинхронно, чтобы не блокировать создание поста
		go func() {
			if err := s.websocketService.BroadcastPost(context.Background(), wsPost); err != nil {
				// Логируем ошибку, но не прерываем создание поста
				// В продакшене здесь можно добавить метрики и алерты
			}
		}()
	}

	return id, nil
}
