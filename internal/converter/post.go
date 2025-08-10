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

func ToPostsFromService(posts []*model.Post) []*api.Post {
	result := make([]*api.Post, 0, len(posts))

	for _, post := range posts {
		result = append(result, ToPostFromService(post))
	}
	return result
}
