package pg

import (
	"context"
	"otus-project/internal/client/db"
	"otus-project/internal/repository/feed"
	feedModel "otus-project/internal/repository/feed/model"
	"time"

	"github.com/google/uuid"
)

type repository struct {
	db db.Client
}

// NewRepository создает новый репозиторий материализованной ленты
func NewRepository(db db.Client) feed.Repository {
	return &repository{
		db: db,
	}
}

// AddToFeed добавляет пост в материализованную ленту пользователя
func (r *repository) AddToFeed(ctx context.Context, userID, postID, authorID, postText string) error {
	query := `
		INSERT INTO materialized_feeds (id, user_id, post_id, author_id, post_text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, post_id) DO UPDATE SET
			post_text = EXCLUDED.post_text,
			updated_at = EXCLUDED.updated_at
	`

	now := time.Now()
	q := db.Query{
		Name:     "feed_repository.AddToFeed",
		QueryRaw: query,
	}

	_, err := r.db.DB().ExecContext(ctx, q,
		uuid.New().String(),
		userID,
		postID,
		authorID,
		postText,
		now,
		now,
	)

	return err
}

// GetFeed получает материализованную ленту пользователя
func (r *repository) GetFeed(ctx context.Context, userID string, offset, limit int) ([]*feedModel.MaterializedFeed, error) {
	query := `
		SELECT id, user_id, post_id, author_id, post_text, created_at, updated_at
		FROM materialized_feeds
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	q := db.Query{
		Name:     "feed_repository.GetFeed",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []*feedModel.MaterializedFeed
	for rows.Next() {
		feed := &feedModel.MaterializedFeed{}
		err := rows.Scan(
			&feed.ID,
			&feed.UserID,
			&feed.PostID,
			&feed.AuthorID,
			&feed.PostText,
			&feed.CreatedAt,
			&feed.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

// RemoveFromFeed удаляет пост из материализованной ленты пользователя
func (r *repository) RemoveFromFeed(ctx context.Context, userID, postID string) error {
	query := `DELETE FROM materialized_feeds WHERE user_id = $1 AND post_id = $2`
	q := db.Query{
		Name:     "feed_repository.RemoveFromFeed",
		QueryRaw: query,
	}
	_, err := r.db.DB().ExecContext(ctx, q, userID, postID)
	return err
}

// CreateJob создает задание на материализацию ленты
func (r *repository) CreateJob(ctx context.Context, job *feedModel.FeedJob) error {
	query := `
		INSERT INTO feed_jobs (id, user_id, post_id, status, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	q := db.Query{
		Name:     "feed_repository.CreateJob",
		QueryRaw: query,
	}

	_, err := r.db.DB().ExecContext(ctx, q,
		job.ID,
		job.UserID,
		job.PostID,
		job.Status,
		job.Priority,
		now,
		now,
	)

	return err
}

// UpdateJobStatus обновляет статус задания
func (r *repository) UpdateJobStatus(ctx context.Context, jobID, status string, error *string) error {
	query := `
		UPDATE feed_jobs 
		SET status = $1, error = $2, updated_at = $3
		WHERE id = $4
	`

	now := time.Now()
	q := db.Query{
		Name:     "feed_repository.UpdateJobStatus",
		QueryRaw: query,
	}
	_, err := r.db.DB().ExecContext(ctx, q, status, error, now, jobID)
	return err
}

// GetPendingJobs получает задания со статусом pending
func (r *repository) GetPendingJobs(ctx context.Context, limit int) ([]*feedModel.FeedJob, error) {
	query := `
		SELECT id, user_id, post_id, status, priority, created_at, updated_at, error
		FROM feed_jobs
		WHERE status = 'pending'
		ORDER BY priority ASC, created_at ASC
		LIMIT $1
	`

	q := db.Query{
		Name:     "feed_repository.GetPendingJobs",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*feedModel.FeedJob
	for rows.Next() {
		job := &feedModel.FeedJob{}
		err := rows.Scan(
			&job.ID,
			&job.UserID,
			&job.PostID,
			&job.Status,
			&job.Priority,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.Error,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetFriendsOfUser получает список друзей пользователя
func (r *repository) GetFriendsOfUser(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT friend_id
		FROM friends
		WHERE user_id = $1 AND status = 'accepted'
		UNION
		SELECT user_id
		FROM friends
		WHERE friend_id = $1 AND status = 'accepted'
	`

	q := db.Query{
		Name:     "feed_repository.GetFriendsOfUser",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []string
	for rows.Next() {
		var friendID string
		if err := rows.Scan(&friendID); err != nil {
			return nil, err
		}
		friends = append(friends, friendID)
	}

	return friends, nil
}
