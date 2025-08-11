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

func ToPostsFromRepo(posts []*modelRepo.Post) []*model.Post {
	var result []*model.Post
	for _, p := range posts {
		result = append(result, ToPostFromRepo(p))
	}
	return result
}
