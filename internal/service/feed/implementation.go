package feed

import (
	"context"
	"log"
	"math/rand"
	"otus-project/internal/client/queue"
	"otus-project/internal/model"
	"otus-project/internal/repository/feed"
	feedModel "otus-project/internal/repository/feed/model"
	"time"

	"github.com/google/uuid"
)

const (
	// MaxFriendsPerPost максимальное количество друзей для обработки одного поста
	// Защита от "эффекта Леди Гаги" - популярные пользователи не должны перегружать систему
	MaxFriendsPerPost = 100
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
	log.Printf("DEBUG: ScheduleFeedUpdate called - postID: %s, authorID: %s", postID, authorID)

	// Получаем список друзей автора поста
	friends, err := s.feedRepository.GetFriendsOfUser(ctx, authorID)
	if err != nil {
		return err
	}

	log.Printf("DEBUG: Got %d friends for user %s", len(friends), authorID)

	// Защита от "эффекта Леди Гаги" - ограничиваем количество друзей для обработки
	friendsToProcess := s.limitFriendsForProcessing(friends, MaxFriendsPerPost)

	log.Printf("DEBUG: Processing %d friends (limited from %d)", len(friendsToProcess), len(friends))

	// Создаем событие ленты
	event := &model.FeedEvent{
		PostID:       postID,
		AuthorUserID: authorID,
		PostText:     postText,
		CreatedAt:    time.Now(),
		EventType:    "post_created",
	}

	// Для каждого выбранного друга создаем задачу обновления ленты
	for _, friendID := range friendsToProcess {
		// Пропускаем автора поста (он не должен видеть свой пост в ленте друзей)
		if friendID == authorID {
			continue
		}

		// Проверяем, что friendID не пустой
		if friendID == "" {
			log.Printf("WARNING: Empty friendID found, skipping")
			continue
		}

		// Определяем приоритет на основе активности пользователя
		priority := s.calculateTaskPriority(authorID, friendID)

		task := &model.FeedUpdateTask{
			UserID:    friendID,
			PostID:    postID,
			Event:     event,
			Priority:  priority,
			CreatedAt: time.Now(),
		}

		log.Printf("DEBUG: Creating task for friend %s, post %s, task.UserID=%s, priority=%d", friendID, postID, task.UserID, priority)

		// Публикуем задачу в очередь
		if err := s.queueClient.PublishFeedUpdateTask(ctx, task); err != nil {
			log.Printf("Error publishing feed update task for user %s: %v", friendID, err)
			continue
		}

		log.Printf("DEBUG: Successfully published task for friend %s", friendID)

		// Отправляем событие через WebSocket для конкретного пользователя
		if err := s.queueClient.PublishFeedEvent(ctx, friendID, event); err != nil {
			log.Printf("Error publishing feed event for user %s: %v", friendID, err)
		}
	}

	log.Printf("Scheduled feed updates for %d/%d friends of user %s (limited to prevent celebrity effect)",
		len(friendsToProcess), len(friends), authorID)
	return nil
}

// calculateTaskPriority вычисляет приоритет задачи на основе активности пользователей
// Высокий приоритет для активных пользователей, низкий для неактивных
func (s *service) calculateTaskPriority(authorID, friendID string) int {
	// Базовая логика приоритизации:
	// 1 - Высокий приоритет (активные пользователи)
	// 3 - Средний приоритет (обычные пользователи)
	// 5 - Низкий приоритет (неактивные пользователи)

	// TODO: В будущем можно добавить более сложную логику:
	// - Анализ частоты постов автора
	// - История взаимодействий между пользователями
	// - Время последней активности друга

	// Пока используем случайный приоритет для демонстрации
	rand.Seed(time.Now().UnixNano() + int64(len(authorID)+len(friendID)))
	randomValue := rand.Intn(100)

	if randomValue < 20 {
		return 1 // Высокий приоритет для 20% задач
	} else if randomValue < 80 {
		return 3 // Средний приоритет для 60% задач
	} else {
		return 5 // Низкий приоритет для 20% задач
	}
}

// limitFriendsForProcessing ограничивает количество друзей для обработки
// Защита от "эффекта Леди Гаги" - если у пользователя слишком много друзей,
// выбираем случайную выборку для обработки
func (s *service) limitFriendsForProcessing(friends []string, maxCount int) []string {
	if len(friends) <= maxCount {
		return friends
	}

	// Создаем копию слайса для перемешивания
	friendsCopy := make([]string, len(friends))
	copy(friendsCopy, friends)

	// Перемешиваем друзей для случайной выборки
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(friendsCopy), func(i, j int) {
		friendsCopy[i], friendsCopy[j] = friendsCopy[j], friendsCopy[i]
	})

	// Возвращаем только первые maxCount друзей
	return friendsCopy[:maxCount]
}

// ProcessFeedUpdateTask обрабатывает задачу обновления ленты
func (s *service) ProcessFeedUpdateTask(ctx context.Context, task *model.FeedUpdateTask) error {
	log.Printf("DEBUG: Processing task - UserID: '%s', PostID: '%s', Priority: %d", task.UserID, task.PostID, task.Priority)

	// Проверяем валидность задачи
	if task.UserID == "" || task.PostID == "" {
		log.Printf("ERROR: Invalid task - UserID: '%s', PostID: '%s'", task.UserID, task.PostID)
		return nil // Пропускаем пустые задачи
	}

	// Проверяем, не обрабатывали ли мы уже эту задачу
	// Это поможет избежать дублирующей обработки
	jobKey := task.UserID + ":" + task.PostID
	log.Printf("DEBUG: Processing job with key: %s", jobKey)

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

	log.Printf("DEBUG: Created job %s for user %s, post %s", job.ID, task.UserID, task.PostID)

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
