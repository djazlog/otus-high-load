package converter

import (
	"otus-project/internal/model"
	"otus-project/pkg/api"
)

func ToPostFromService(post *model.Post) *api.Post {

	return &api.Post{
		Id:           post.ID,
		Text:         post.Text,
		AuthorUserId: post.AuthorUserId,
	}
}
