package dialog

import (
	"context"
	"otus-project/internal/model"
)

// DialogService интерфейс сервиса диалогов
type DialogService interface {
	// SendMessage отправляет сообщение в диалог
	SendMessage(ctx context.Context, fromUserId, toUserId string, text string) error
	// GetDialogList возвращает список сообщений диалога между двумя пользователями
	GetDialogList(ctx context.Context, userId1, userId2 string) ([]*model.DialogMessage, error)
	// GetDialogCount возвращает общее количество диалогов (только для Redis)
	GetDialogCount(ctx context.Context) (int64, error)
	// GetDialogStats возвращает статистику диалогов (только для Redis)
	GetDialogStats(ctx context.Context) (totalDialogs, activeDialogs int64, err error)
}
