package converter

import (
	"otus-project/internal/model"
	modelRepo "otus-project/internal/repository/post/redis/model"
)

func ToPostFromRepo(post *modelRepo.Post) *model.Post {
	return &model.Post{
		ID:           &post.ID,
		Text:         &post.Text,
		AuthorUserId: &post.AuthorUserId,
	}
}
