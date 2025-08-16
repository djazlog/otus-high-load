package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"otus-project/internal/client/cache"
	"otus-project/internal/model"
	"otus-project/internal/repository"
	"otus-project/internal/repository/post/redis/converter"
	modelRepo "otus-project/internal/repository/post/redis/model"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var (
	keyList  = "posts"
	redisTTL = time.Minute * 1 // Время жизни кэша
)

type repo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) repository.PostRepository {
	return &repo{cl: cl}
}

func (r *repo) Create(ctx context.Context, post *model.Post) (*string, error) {
	id := uuid.New().String()

	newPost := modelRepo.Post{
		ID:           id,
		Text:         *post.Text,
		AuthorUserId: *post.AuthorUserId,
		CreatedAtNs:  time.Now().UnixNano(),
	}

	err := r.cl.HashSet(ctx, keyList, newPost, redisTTL)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *repo) Get(ctx context.Context, offset *float32, limit *float32) (*model.Post, error) {
	values, err := r.cl.HGetAll(ctx, keyList)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorPostNotFound
	}

	var post modelRepo.Post
	err = redigo.ScanStruct(values, &post)
	if err != nil {
		return nil, err
	}

	return converter.ToPostFromRepo(&post), nil
}

// Feed Получем посты друзей пользователя
func (r *repo) Feed(ctx context.Context, id string, offset *float32, limit *float32) ([]*model.Post, error) {
	cacheKey := fmt.Sprintf("feed:user:%s", id)

	// Получаем JSON-строку из Redis
	value, err := r.cl.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, model.ErrorPostNotFound
	}

	// Проверяем тип значения и конвертируем его в string
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	default:
		return nil, fmt.Errorf("unexpected type for cached value: %T", value)
	}

	var cachedPosts []*model.Post
	err = json.Unmarshal([]byte(strValue), &cachedPosts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal posts from cache: %w", err)
	}

	return cachedPosts, nil
}

// CacheFeed сохраняем посты в кэш
func (r *repo) CacheFeed(ctx context.Context, userId string, posts []*model.Post) error {
	// Формируем ключ кэша
	cacheKey := fmt.Sprintf("feed:user:%s", userId)

	// Сериализуем посты в JSON
	postJSON, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("failed to marshal posts to JSON: %w", err)
	}

	// Сохраняем данные в Redis с TTL
	err = r.cl.Set(ctx, cacheKey, postJSON, redisTTL)
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// GetByID получает пост по ID
func (r *repo) GetByID(ctx context.Context, id string) (*model.Post, error) {
	// В Redis репозитории GetByID не реализован, так как это кэш
	// Возвращаем ошибку, чтобы использовать PostgreSQL
	return nil, fmt.Errorf("GetByID not implemented in Redis repository")
}

// Update обновляет пост
func (r *repo) Update(ctx context.Context, id string, text string) error {
	// В Redis репозитории Update не реализован, так как это кэш
	// Возвращаем ошибку, чтобы использовать PostgreSQL
	return fmt.Errorf("Update not implemented in Redis repository")
}

// Delete удаляет пост
func (r *repo) Delete(ctx context.Context, id string) error {
	// В Redis репозитории Delete не реализован, так как это кэш
	// Возвращаем ошибку, чтобы использовать PostgreSQL
	return fmt.Errorf("Delete not implemented in Redis repository")
}
