package friend

import (
	"otus-project/internal/client/db"

	"otus-project/internal/repository"
	"otus-project/internal/service"
)

type serv struct {
	friendRepository repository.FriendRepository
	txManager        db.TxManager
}

func NewService(
	friendRepository repository.FriendRepository,
	txManager db.TxManager,
) service.FriendService {
	return &serv{
		friendRepository: friendRepository,
		txManager:        txManager,
	}
}
