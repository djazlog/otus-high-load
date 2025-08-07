package converter

import (
	"otus-project/internal/model"
	modelRepo "otus-project/internal/repository/post/pg/model"
)

func ToPostFromRepo(post *modelRepo.Post) *model.Post {
	return &model.Post{
		ID:           post.ID,
		Text:         post.Text,
		AuthorUserId: post.AuthorUserId,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}
}
