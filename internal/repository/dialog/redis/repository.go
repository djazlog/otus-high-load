package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"otus-project/internal/client/cache"
	"otus-project/internal/model"
	"otus-project/internal/repository"
	"otus-project/internal/repository/dialog/redis/converter"
	repoModel "otus-project/internal/repository/dialog/redis/model"
	"otus-project/internal/utils"
	"time"

	"github.com/google/uuid"
)

var (
	dialogKeyPrefix = "dialog"
	redisTTL        = time.Hour * 24 // Время жизни кэша диалогов
)

type repo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) repository.DialogRepository {
	return &repo{cl: cl}
}

// SendMessage сохраняет сообщение в диалоге в Redis с использованием LUA скрипта
func (r *repo) SendMessage(ctx context.Context, fromUserId, toUserId, text string) error {
	dialogKey := utils.GenerateDialogKey(fromUserId, toUserId)
	messageID := uuid.New().String()

	// Создаем сообщение
	message := repoModel.DialogMessage{
		ID:          messageID,
		FromUserID:  fromUserId,
		ToUserID:    toUserId,
		Text:        text,
		CreatedAtNs: time.Now().UnixNano(),
		DialogKey:   dialogKey,
	}

	// Сериализуем сообщение в JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	// Ключи для LUA скрипта
	dialogRedisKey := fmt.Sprintf("%s:%s", dialogKeyPrefix, dialogKey)
	counterKey := "dialog:counter:total"

	// Выполняем LUA скрипт для сохранения сообщения и увеличения счетчика
	_, err = r.cl.Eval(ctx, SendMessageScript, []string{dialogRedisKey, counterKey}, messageJSON, int(redisTTL.Seconds()))
	if err != nil {
		return fmt.Errorf("failed to execute LUA script for sending message: %w", err)
	}

	return nil
}

// GetDialogList возвращает список сообщений диалога между двумя пользователями
func (r *repo) GetDialogList(ctx context.Context, userId1, userId2 string) ([]*model.DialogMessage, error) {
	dialogKey := utils.GenerateDialogKey(userId1, userId2)
	redisKey := fmt.Sprintf("%s:%s", dialogKeyPrefix, dialogKey)

	// Получаем все сообщения из Redis списка
	messagesJSON, err := r.cl.LRange(ctx, redisKey, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from Redis: %w", err)
	}

	if len(messagesJSON) == 0 {
		return []*model.DialogMessage{}, nil
	}

	var messages []*repoModel.DialogMessage
	for _, msgJSON := range messagesJSON {
		var message repoModel.DialogMessage

		// Проверяем тип значения и конвертируем его в string
		var strValue string
		switch v := msgJSON.(type) {
		case string:
			strValue = v
		case []byte:
			strValue = string(v)
		default:
			return nil, fmt.Errorf("unexpected type for cached message: %T", msgJSON)
		}

		err := json.Unmarshal([]byte(strValue), &message)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal message from JSON: %w", err)
		}

		messages = append(messages, &message)
	}

	// Конвертируем в сервисные модели
	return converter.ToDialogMessagesFromRepo(messages), nil
}

// GetDialogCount возвращает общее количество диалогов в Redis
func (r *repo) GetDialogCount(ctx context.Context) (int64, error) {
	counterKey := "dialog:counter:total"

	result, err := r.cl.Eval(ctx, GetDialogCountScript, []string{counterKey})
	if err != nil {
		return 0, fmt.Errorf("failed to execute LUA script for getting dialog count: %w", err)
	}

	// Конвертируем результат в int64
	switch v := result.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unexpected result type for dialog count: %T", result)
	}
}

// GetDialogStats возвращает статистику диалогов
func (r *repo) GetDialogStats(ctx context.Context) (totalDialogs, activeDialogs int64, err error) {
	counterKey := "dialog:counter:total"
	dialogPattern := fmt.Sprintf("%s:*", dialogKeyPrefix)

	result, err := r.cl.Eval(ctx, GetDialogStatsScript, []string{counterKey}, dialogPattern)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute LUA script for getting dialog stats: %w", err)
	}

	// Результат должен быть массивом из двух элементов
	results, ok := result.([]interface{})
	if !ok || len(results) != 2 {
		return 0, 0, fmt.Errorf("unexpected result format for dialog stats")
	}

	// Конвертируем первый элемент (total dialogs)
	switch v := results[0].(type) {
	case int64:
		totalDialogs = v
	case int:
		totalDialogs = int64(v)
	case float64:
		totalDialogs = int64(v)
	default:
		return 0, 0, fmt.Errorf("unexpected type for total dialogs: %T", results[0])
	}

	// Конвертируем второй элемент (active dialogs)
	switch v := results[1].(type) {
	case int64:
		activeDialogs = v
	case int:
		activeDialogs = int64(v)
	case float64:
		activeDialogs = int64(v)
	default:
		return 0, 0, fmt.Errorf("unexpected type for active dialogs: %T", results[1])
	}

	return totalDialogs, activeDialogs, nil
}
