package converter

import (
	"otus-project/internal/model"
	"otus-project/pkg/api"
)

// ToDialogMessageFromService конвертирует модель диалога в API модель
func ToDialogMessageFromService(msg *model.DialogMessage) *api.DialogMessage {
	return &api.DialogMessage{
		From: api.UserId(msg.From),
		To:   api.UserId(msg.To),
		Text: api.DialogMessageText(msg.Text),
	}
}

// ToDialogMessagesFromService конвертирует список моделей диалогов в API модели
func ToDialogMessagesFromService(messages []*model.DialogMessage) []api.DialogMessage {
	var result []api.DialogMessage
	for _, msg := range messages {
		if msg == nil {
			continue
		}
		result = append(result, *ToDialogMessageFromService(msg))
	}
	if len(result) == 0 {
		result = []api.DialogMessage{}
	}
	return result
}
