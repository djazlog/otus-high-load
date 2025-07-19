package converter

import (
	"otus-project/internal/model"
	"otus-project/pkg/api"
)

func ToUserFromService(user *model.UserInfo) *api.User {
	var birthday api.BirthDate
	birthday.Time = *user.Birthdate

	return &api.User{

		Id:         user.Id,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		City:       user.City,
		Birthdate:  &birthday,
		Biography:  user.Biography,
	}
}

func ToUsersFromService(users []*model.UserInfo) []api.User {
	var result []api.User
	for _, u := range users {
		if u == nil {
			continue
		}
		result = append(result, *ToUserFromService(u))
	}
	if len(result) == 0 {
		result = []api.User{}
	}
	return result
}

func ToUserInfoFromApi(info *api.PostUserRegisterJSONBody) *model.UserInfo {
	return &model.UserInfo{
		FirstName:  info.FirstName,
		SecondName: info.SecondName,
		City:       info.City,
		Birthdate:  &info.Birthdate.Time,
		Biography:  info.Biography,
		Password:   info.Password,
	}
}

func ToUserFilterFromApi(info *api.GetUserSearchParams) *model.UserFilter {
	return &model.UserFilter{
		FirstName: info.FirstName,
		LastName:  info.LastName,
	}
}

type TokenJSONBody struct {
	// Text Текст поста
	Token string `json:"token"`
}

func ToTokenResponse(obj *string) *TokenJSONBody {
	return &TokenJSONBody{
		Token: *obj,
	}
}
