package redis

import (
	"context"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"otus-project/internal/client/cache"
	"otus-project/internal/model"
	"otus-project/internal/repository"
	"otus-project/internal/repository/post/redis/converter"
	modelRepo "otus-project/internal/repository/post/redis/model"
	"time"
)

var (
	keyList = "posts"
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

	err := r.cl.HashSet(ctx, keyList, newPost)
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
