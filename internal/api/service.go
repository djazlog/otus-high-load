package api

import (
	"otus-project/internal/service"
)

type Implementation struct {
	userService   service.UserService
	postService   service.PostService
	friendService service.FriendService
	dialogService service.DialogService
	feedService   service.FeedService
}

func NewImplementation(
	userService service.UserService,
	postService service.PostService,
	friendService service.FriendService,
	dialogService service.DialogService,
) *Implementation {
	return &Implementation{
		userService:   userService,
		postService:   postService,
		friendService: friendService,
		dialogService: dialogService,
	}
}
