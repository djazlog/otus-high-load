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

type TokenJSONBody struct {
	// Text Текст поста
	Token string `json:"token"`
}

func ToTokenResponse(obj *string) *TokenJSONBody {
	return &TokenJSONBody{
		Token: *obj,
	}
}
