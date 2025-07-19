package repository

import (
	"context"
	"otus-project/internal/model"
)

type UserRepository interface {
	Login(ctx context.Context, login *model.LoginDto) (*string, error)
	Register(ctx context.Context, info *model.UserInfo) (string, error)
	Get(ctx context.Context, id string) (*model.UserInfo, error)
	Search(ctx context.Context, filter *model.UserFilter) ([]*model.UserInfo, error)
}
