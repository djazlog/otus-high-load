package feed

import (
	"context"
	"log"
	"otus-project/internal/client/queue"
	"otus-project/internal/model"
	"otus-project/internal/repository/feed"
	feedModel "otus-project/internal/repository/feed/model"
	"time"

	"github.com/google/uuid"
)

type service struct {
	feedRepository feed.Repository
	queueClient    queue.Client
	workerCtx      context.Context
	workerCancel   context.CancelFunc
}

// NewService создает новый сервис отложенной материализации ленты
func NewService(feedRepository feed.Repository, queueClient queue.Client) Service {
	return &service{
		feedRepository: feedRepository,
		queueClient:    queueClient,
	}
}

// ScheduleFeedUpdate планирует обновление ленты для друзей автора поста
func (s *service) ScheduleFeedUpdate(ctx context.Context, postID, authorID, postText string) error {
	// Получаем список друзей автора поста
	friends, err := s.feedRepository.GetFriendsOfUser(ctx, authorID)
	if err != nil {
		return err
	}

	// Создаем событие ленты
	event := &model.FeedEvent{
		PostID:       postID,
		AuthorUserID: authorID,
		PostText:     postText,
		CreatedAt:    time.Now(),
		EventType:    "post_created",
	}

	// Для каждого друга создаем задачу обновления ленты
	for _, friendID := range friends {
		// Пропускаем автора поста (он не должен видеть свой пост в ленте друзей)
		if friendID == authorID {
			continue
		}

		task := &model.FeedUpdateTask{
			UserID:    friendID,
			PostID:    postID,
			Event:     event,
			Priority:  3, // Средний приоритет
			CreatedAt: time.Now(),
		}

		// Публикуем задачу в очередь
		if err := s.queueClient.PublishFeedUpdateTask(ctx, task); err != nil {
			log.Printf("Error publishing feed update task for user %s: %v", friendID, err)
			continue
		}

		// Отправляем событие через WebSocket для конкретного пользователя
		if err := s.queueClient.PublishFeedEvent(ctx, friendID, event); err != nil {
			log.Printf("Error publishing feed event for user %s: %v", friendID, err)
		}
	}

	log.Printf("Scheduled feed updates for %d friends of user %s", len(friends), authorID)
	return nil
}

// ProcessFeedUpdateTask обрабатывает задачу обновления ленты
func (s *service) ProcessFeedUpdateTask(ctx context.Context, task *model.FeedUpdateTask) error {
	// Создаем задание в БД для отслеживания
	job := &feedModel.FeedJob{
		ID:        uuid.New().String(),
		UserID:    task.UserID,
		PostID:    task.PostID,
		Status:    "processing",
		Priority:  task.Priority,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.feedRepository.CreateJob(ctx, job); err != nil {
		log.Printf("Error creating feed job: %v", err)
		return err
	}

	// Добавляем пост в материализованную ленту пользователя
	if err := s.feedRepository.AddToFeed(ctx, task.UserID, task.PostID, task.Event.AuthorUserID, task.Event.PostText); err != nil {
		// Обновляем статус задания на failed
		errorMsg := err.Error()
		if updateErr := s.feedRepository.UpdateJobStatus(ctx, job.ID, "failed", &errorMsg); updateErr != nil {
			log.Printf("Error updating job status: %v", updateErr)
		}
		return err
	}

	// Обновляем статус задания на completed
	if err := s.feedRepository.UpdateJobStatus(ctx, job.ID, "completed", nil); err != nil {
		log.Printf("Error updating job status: %v", err)
	}

	log.Printf("Processed feed update task for user %s, post %s", task.UserID, task.PostID)
	return nil
}

// GetMaterializedFeed получает материализованную ленту пользователя
func (s *service) GetMaterializedFeed(ctx context.Context, userID string, offset, limit int) ([]*feedModel.MaterializedFeed, error) {
	return s.feedRepository.GetFeed(ctx, userID, offset, limit)
}

// StartWorker запускает воркер для обработки задач материализации
func (s *service) StartWorker(ctx context.Context) error {
	s.workerCtx, s.workerCancel = context.WithCancel(ctx)

	// Запускаем потребление задач из очереди
	err := s.queueClient.ConsumeFeedMaterializationTasks(s.workerCtx, s.ProcessFeedUpdateTask)
	if err != nil {
		return err
	}

	log.Println("Feed materialization worker started")
	return nil
}

// StopWorker останавливает воркер
func (s *service) StopWorker(ctx context.Context) error {
	if s.workerCancel != nil {
		s.workerCancel()
		log.Println("Feed materialization worker stopped")
	}
	return nil
}
