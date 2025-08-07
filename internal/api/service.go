package api

import (
	"net/http"
	"otus-project/internal/service"
	"otus-project/pkg/api"
)

type Implementation struct {
	userService   service.UserService
	postService   service.PostService
	friendService service.FriendService
}

func NewImplementation(
	userService service.UserService,
	postService service.PostService,
	friendService service.FriendService,
) *Implementation {
	return &Implementation{
		userService:   userService,
		postService:   postService,
		friendService: friendService,
	}
}

func (i *Implementation) GetDialogUserIdList(w http.ResponseWriter, r *http.Request, userId api.UserId) {
	//TODO implement me
	panic("implement me")
}

func (i *Implementation) PostDialogUserIdSend(w http.ResponseWriter, r *http.Request, userId api.UserId) {
	//TODO implement me
	panic("implement me")
}

func (i *Implementation) PutFriendDeleteUserId(w http.ResponseWriter, r *http.Request, userId api.UserId) {
	//TODO implement me
	panic("implement me")
}

func (i *Implementation) PostPostCreate(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (i *Implementation) PutPostDeleteId(w http.ResponseWriter, r *http.Request, id api.PostId) {
	//TODO implement me
	panic("implement me")
}

func (i *Implementation) GetPostGetId(w http.ResponseWriter, r *http.Request, id api.PostId) {
	//TODO implement me
	panic("implement me")
}

func (i *Implementation) PutPostUpdate(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
