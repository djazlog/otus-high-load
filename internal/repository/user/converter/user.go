package converter

import (
	"otus-project/internal/model"
	modelRepo "otus-project/internal/repository/user/model"
)

// ToUserInfoFromRepo - convert modelRepo.User to model.UserInfo
func ToUserInfoFromRepo(info *modelRepo.User) *model.UserInfo {
	return &model.UserInfo{
		Id:         info.Id,
		FirstName:  info.FirstName,
		SecondName: info.SecondName,
		Biography:  info.Biography,
		Birthdate:  info.Birthdate,
		City:       info.City,
		CreatedAt:  info.CreatedAt,
		UpdatedAt:  info.UpdatedAt,
	}
}
