package dialog

import (
	"context"
	"fmt"
	"otus-project/internal/model"
	"otus-project/internal/repository"
)

type Implementation struct {
	dialogPgRepo    repository.DialogRepository
	dialogRedisRepo repository.DialogRepository
}

func NewImplementation(dialogPgRepo, dialogRedisRepo repository.DialogRepository) *Implementation {
	return &Implementation{
		dialogPgRepo:    dialogPgRepo,
		dialogRedisRepo: dialogRedisRepo,
	}
}

// SendMessage сохраняет сообщение в диалоге
func (i *Implementation) SendMessage(ctx context.Context, fromUserId, toUserId string, text string) error {
	// Сначала пытаемся сохранить в Redis
	err := i.dialogRedisRepo.SendMessage(ctx, fromUserId, toUserId, text)
	if err != nil {
		// Если не удалось, сохраняем в PostgreSQL
		//return i.dialogPgRepo.SendMessage(ctx, fromUserId, toUserId, text)
		return err
	}

	return nil
}

// GetDialogList получает список сообщений диалога между двумя пользователями
func (i *Implementation) GetDialogList(ctx context.Context, userId1, userId2 string) ([]*model.DialogMessage, error) {
	// Сначала пытаемся получить из Redis
	messages, err := i.dialogRedisRepo.GetDialogList(ctx, userId1, userId2)
	if err != nil {
		fmt.Printf("Failed to get dialog from Redis: %v, trying PostgreSQL\n", err)
		// Если не удалось, получаем из PostgreSQL
		//return i.dialogPgRepo.GetDialogList(ctx, userId1, userId2)
		return nil, err
	}

	return messages, nil

	/*if len(messages) > 0 {
		fmt.Println("dialog messages from redis")
		return messages, nil
	}

	// Если в Redis нет сообщений, получаем из PostgreSQL
	messages, err = i.dialogPgRepo.GetDialogList(ctx, userId1, userId2)
	if err != nil {
		return nil, err
	}

	fmt.Println("dialog messages from postgres")
	return messages, nil*/
}

// GetDialogCount возвращает общее количество диалогов
func (i *Implementation) GetDialogCount(ctx context.Context) (int64, error) {
	// Пытаемся получить из Redis
	count, err := i.dialogRedisRepo.GetDialogCount(ctx)
	if err != nil {
		fmt.Printf("Failed to get dialog count from Redis: %v\n", err)
		return 0, err
	}

	return count, nil
}

// GetDialogStats возвращает статистику диалогов
func (i *Implementation) GetDialogStats(ctx context.Context) (totalDialogs, activeDialogs int64, err error) {
	// Пытаемся получить из Redis
	total, active, err := i.dialogRedisRepo.GetDialogStats(ctx)
	if err != nil {
		fmt.Printf("Failed to get dialog stats from Redis: %v\n", err)
		return 0, 0, err
	}

	return total, active, nil
}
