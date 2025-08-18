package app

import (
	"context"
	"log"
	"otus-project/internal/api"
	"otus-project/internal/client/cache/redis"

	"otus-project/internal/client/cache"
	"otus-project/internal/client/db"
	"otus-project/internal/client/db/pg"
	"otus-project/internal/client/db/transaction"
	"otus-project/internal/client/queue"
	"otus-project/internal/client/queue/rabbitmq"
	"otus-project/internal/closer"
	"otus-project/internal/config"
	"otus-project/internal/repository"
	dialogRepo "otus-project/internal/repository/dialog"
	feedRepo "otus-project/internal/repository/feed"
	feedPgRepo "otus-project/internal/repository/feed/pg"
	friendRepo "otus-project/internal/repository/friend"
	postPgRepo "otus-project/internal/repository/post/pg"
	postRRepo "otus-project/internal/repository/post/redis"
	userRepository "otus-project/internal/repository/user"
	"otus-project/internal/service"
	dialogService "otus-project/internal/service/dialog"
	eventBusService "otus-project/internal/service/event_bus"
	feedService "otus-project/internal/service/feed"
	friendService "otus-project/internal/service/friend"
	postService "otus-project/internal/service/post"
	userService "otus-project/internal/service/user"
	websocketService "otus-project/internal/service/websocket"

	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	httpConfig  config.HTTPConfig
	redisConfig config.RedisConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	userRepository      repository.UserRepository
	postPgRepository    repository.PostRepository
	postRedisRepository repository.PostRepository
	friendRepository    repository.FriendRepository
	dialogRepository    repository.DialogRepository

	userService      service.UserService
	postService      service.PostService
	friendService    service.FriendService
	dialogService    service.DialogService
	websocketService websocketService.WebSocketService
	feedService      feedService.Service
	queueClient      queue.Client
	eventBus         eventBusService.EventBus

	apiImpl *api.Implementation
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

// RedisConfig возвращает конфиг redis
func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

// RedisPool возвращает пул соединений к redis
func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

// RedisClient возвращает клиент redis
func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
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

// PostRepository возвращает репозиторий Post
func (s *serviceProvider) PostRepository(ctx context.Context) repository.PostRepository {
	if s.postPgRepository == nil {
		s.postPgRepository = postPgRepo.NewRepository(s.DBClient(ctx))
	}

	return s.postPgRepository
}

// PostRedisRepository возвращает репозиторий Post
func (s *serviceProvider) PostRedisRepository(ctx context.Context) repository.PostRepository {
	if s.postRedisRepository == nil {
		s.postRedisRepository = postRRepo.NewRepository(s.RedisClient())
	}

	return s.postRedisRepository
}

// FriendRepository возвращает репозиторий Post
func (s *serviceProvider) FriendRepository(ctx context.Context) repository.FriendRepository {
	if s.friendRepository == nil {
		s.friendRepository = friendRepo.NewRepository(s.DBClient(ctx))
	}

	return s.friendRepository
}

// DialogRepository возвращает репозиторий диалогов
func (s *serviceProvider) DialogRepository(ctx context.Context) repository.DialogRepository {
	if s.dialogRepository == nil {
		s.dialogRepository = dialogRepo.NewRepository(s.DBClient(ctx))
	}

	return s.dialogRepository
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

// PostService возвращает сервис
func (s *serviceProvider) PostService(ctx context.Context) service.PostService {
	if s.postService == nil {
		s.postService = postService.NewService(
			s.PostRepository(ctx),
			s.PostRedisRepository(ctx),
			s.TxManager(ctx),
			s.EventBus(),
		)
	}

	return s.postService
}

// FriendService возвращает сервис User
func (s *serviceProvider) FriendService(ctx context.Context) service.FriendService {
	if s.friendService == nil {
		s.friendService = friendService.NewService(
			s.FriendRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.friendService
}

// DialogService возвращает сервис диалогов
func (s *serviceProvider) DialogService(ctx context.Context) service.DialogService {
	if s.dialogService == nil {
		s.dialogService = dialogService.NewImplementation(s.DialogRepository(ctx))
	}

	return s.dialogService
}

// WebSocketService возвращает WebSocket сервис
func (s *serviceProvider) WebSocketService() websocketService.WebSocketService {
	if s.websocketService == nil {
		s.websocketService = websocketService.NewService()
	}

	return s.websocketService
}

// QueueClient возвращает клиент очереди сообщений
func (s *serviceProvider) QueueClient() queue.Client {
	if s.queueClient == nil {
		cfg, err := config.NewRabbitMQConfig()
		if err != nil {
			log.Fatalf("failed to get rabbitmq config: %s", err.Error())
		}

		client, err := rabbitmq.NewClient(cfg)
		if err != nil {
			log.Fatalf("failed to create queue client: %s", err.Error())
		}

		s.queueClient = client
		closer.Add(client.Close)
	}

	return s.queueClient
}

// FeedRepository возвращает репозиторий материализованной ленты
func (s *serviceProvider) FeedRepository(ctx context.Context) feedRepo.Repository {
	return feedPgRepo.NewRepository(s.DBClient(ctx))
}

// EventBus возвращает Event Bus
func (s *serviceProvider) EventBus() eventBusService.EventBus {
	if s.eventBus == nil {
		s.eventBus = eventBusService.NewService()
	}

	return s.eventBus
}

// FeedService возвращает сервис отложенной материализации ленты
func (s *serviceProvider) FeedService(ctx context.Context) feedService.Service {
	if s.feedService == nil {
		s.feedService = feedService.NewService(s.FeedRepository(ctx), s.QueueClient())
	}

	return s.feedService
}

// ApiImpl возвращает реализацию сервиса User
func (s *serviceProvider) ApiImpl(ctx context.Context) *api.Implementation {
	if s.apiImpl == nil {
		s.apiImpl = api.NewImplementation(s.UserService(ctx), s.PostService(ctx), s.FriendService(ctx), s.DialogService(ctx))
	}

	return s.apiImpl
}
