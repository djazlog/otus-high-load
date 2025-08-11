package converter

import (
	"otus-project/internal/model"
	modelRepo "otus-project/internal/repository/friend/model"
)

// ToFriendFromRepo - convert modelRepo.Friend to model.Friend
func ToFriendFromRepo(info *modelRepo.Friend) *model.Friend {
	return &model.Friend{
		FriendId:  info.FriendId,
		UserId:    info.UserId,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}
