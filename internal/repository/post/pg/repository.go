package pg

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"otus-project/internal/client/db"
	"otus-project/internal/model"
	"otus-project/internal/repository"
	"otus-project/internal/repository/post/pg/converter"
	modelRepo "otus-project/internal/repository/post/pg/model"
)

const (
	tableName = "posts"

	idColumn           = "id"
	textColumn         = "content"
	authorUserIdColumn = "author_user_id"
	userIdColumn       = "user_id"
	createdAtColumn    = "created_at"
	updatedAtColumn    = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.PostRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, post *model.Post) (*string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(idColumn, textColumn, authorUserIdColumn).
		Values(uid, post.Text, post.AuthorUserId).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "post_repository.Create",
		QueryRaw: query,
	}

	var id string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *repo) Get(ctx context.Context, offset *float32, limit *float32) (*model.Post, error) {
	if offset == nil {
		val := float32(1)
		offset = &val
	}
	if limit == nil {
		val := float32(100)
		limit = &val
	}
	// Приводим float32 к uint64
	off := uint64(*offset)
	lim := uint64(*limit)
	builder := sq.Select(idColumn, textColumn, authorUserIdColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Offset(off).
		Limit(lim)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "post_repository.Get",
		QueryRaw: query,
	}

	var post modelRepo.Post
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&post.ID, &post.Text, &post.AuthorUserId, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrorPostNotFound
		}

		return nil, err
	}

	return converter.ToPostFromRepo(&post), nil
}

// Feed Получить посты друзей
func (r *repo) Feed(ctx context.Context, id string, offset *float32, limit *float32) ([]*model.Post, error) {
	if offset == nil {
		val := float32(1)
		offset = &val
	}
	if limit == nil {
		val := float32(1000)
		limit = &val
	}
	// Приводим float32 к uint64
	off := uint64(*offset)
	lim := uint64(*limit)

	builder := sq.Select(idColumn, textColumn, authorUserIdColumn, "posts.created_at", "posts.updated_at").
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{userIdColumn: id}).
		Join("friends f ON f.friend_id = posts.author_user_id").
		OrderBy("posts.created_at DESC").
		Offset(off).
		Limit(lim)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "post_repository.Feed",
		QueryRaw: query,
	}

	rows, err := r.db.ReplicaDB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p modelRepo.Post
		if err := rows.Scan(&p.ID, &p.Text, &p.AuthorUserId, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, converter.ToPostFromRepo(&p))
	}

	return posts, nil
}

// CacheFeed Кэш постов друзей для редиса. Тут как болванка
func (r *repo) CacheFeed(ctx context.Context, userId string, posts []*model.Post) error {
	return nil
}
