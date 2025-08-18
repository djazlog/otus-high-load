package converter

import (
	"otus-project/internal/model"
	repoModel "otus-project/internal/repository/dialog/model"
)

// ToDialogMessageFromRepo конвертирует модель репозитория в сервисную модель
func ToDialogMessageFromRepo(msg *repoModel.DialogMessage) *model.DialogMessage {
	return &model.DialogMessage{
		From:      msg.FromUserID,
		To:        msg.ToUserID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

// ToDialogMessagesFromRepo конвертирует список моделей репозитория в сервисные модели
func ToDialogMessagesFromRepo(messages []*repoModel.DialogMessage) []*model.DialogMessage {
	var result []*model.DialogMessage
	for _, msg := range messages {
		if msg == nil {
			continue
		}
		result = append(result, ToDialogMessageFromRepo(msg))
	}
	if len(result) == 0 {
		result = []*model.DialogMessage{}
	}
	return result
}
