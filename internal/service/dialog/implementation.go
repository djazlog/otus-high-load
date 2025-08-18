package dialog

import (
	"context"
	"otus-project/internal/model"
	"otus-project/internal/repository"
)

type Implementation struct {
	dialogRepo repository.DialogRepository
}

func NewImplementation(dialogRepo repository.DialogRepository) *Implementation {
	return &Implementation{
		dialogRepo: dialogRepo,
	}
}

func (i *Implementation) SendMessage(ctx context.Context, fromUserId, toUserId string, text string) error {
	return i.dialogRepo.SendMessage(ctx, fromUserId, toUserId, text)
}

func (i *Implementation) GetDialogList(ctx context.Context, userId1, userId2 string) ([]*model.DialogMessage, error) {
	return i.dialogRepo.GetDialogList(ctx, userId1, userId2)
}
