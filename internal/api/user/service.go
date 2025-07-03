package user

import (
	"net/http"
	"otus-project/internal/service"
	"otus-project/pkg/api"
)

type Implementation struct {
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
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

func (i *Implementation) PutFriendSetUserId(w http.ResponseWriter, r *http.Request, userId api.UserId) {
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

func (i *Implementation) GetPostFeed(w http.ResponseWriter, r *http.Request, params api.GetPostFeedParams) {
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

func (i *Implementation) GetUserSearch(w http.ResponseWriter, r *http.Request, params api.GetUserSearchParams) {
	//TODO implement me
	panic("implement me")
}
