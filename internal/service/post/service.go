package post

import (
	"otus-project/internal/client/db"
	"otus-project/internal/repository"
	"otus-project/internal/service"
)

type serv struct {
	postPgRepository repository.PostRepository
	postRRepository  repository.PostRepository
	txManager        db.TxManager
}

func NewService(
	postPgRepository repository.PostRepository,
	postRRepository repository.PostRepository,
	txManager db.TxManager,
) service.PostService {
	return &serv{
		postPgRepository: postPgRepository,
		postRRepository:  postRRepository,
		txManager:        txManager,
	}
}
