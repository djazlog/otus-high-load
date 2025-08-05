package app

import (
	"context"
	"log"
	"otus-project/internal/api/user"
	"otus-project/internal/client/db"
	"otus-project/internal/client/db/pg"
	"otus-project/internal/client/db/transaction"
	"otus-project/internal/closer"
	"otus-project/internal/config"
	"otus-project/internal/repository"
	userRepository "otus-project/internal/repository/user"
	"otus-project/internal/service"
	userService "otus-project/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	httpConfig config.HTTPConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository

	userService service.UserService

	userImpl *user.Implementation
}

// NewServiceProvider создает новый сервисный провайдер
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig возвращает конфиг БД
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// HTTPConfig возвращает конфиг http
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

// DBClient возвращает клиент БД
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN(), s.PGConfig().DSNReplica())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager возвращает менеджер транзакций
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// UserRepository возвращает репозиторий User
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// UserService возвращает сервис User
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

// UserImpl возвращает реализацию сервиса User
func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
