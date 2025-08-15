package post

import (
	"context"
	"otus-project/internal/model"
)

type PostRepository interface {
	Create(ctx context.Context, info *model.Post) (*string, error)
	Get(ctx context.Context, offset *float32, limit *float32) (*model.Post, error)
	GetByID(ctx context.Context, id string) (*model.Post, error)
	Update(ctx context.Context, id string, text string) error
	Delete(ctx context.Context, id string) error
	Feed(ctx context.Context, id string, offset *float32, limit *float32) ([]*model.Post, error)
	CacheFeed(ctx context.Context, userId string, posts []*model.Post) error
}
