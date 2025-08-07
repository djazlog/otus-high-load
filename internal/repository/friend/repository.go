package user

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"otus-project/internal/client/db"
	"otus-project/internal/repository"
)

const (
	tableName = "users"

	idColumn     = "user_id"
	friendColumn = "friend_id"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.FriendRepository {
	return &repo{db: db}
}

// AddFriend добавляет нового друга к пользователю.
func (r *repo) AddFriend(ctx context.Context, userId, friendId string) error {
	// Проверяем, существует ли уже связь
	exists, err := r.isFriendExists(ctx, userId, friendId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("пользователь %s уже является другом у %s", friendId, userId)
	}

	// Добавляем новую связь
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(idColumn, friendColumn).
		Values(userId, friendId)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "friend_repository.AddFriend",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("не удалось добавить друга %s для пользователя %s", friendId, userId)
	}

	return nil
}

// isFriendExists проверяет, существует ли связь между пользователем и другом.
func (r *repo) isFriendExists(ctx context.Context, userId, friendId string) (bool, error) {
	builder := sq.Select("COUNT(*)").
		From(tableName).
		Where(sq.Eq{idColumn: userId, friendColumn: friendId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "friend_repository.isFriendExists",
		QueryRaw: query,
	}

	var count int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete удаляет связь между пользователем и другом по их ID.
func (r *repo) Delete(ctx context.Context, userId, friendId string) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userId, friendColumn: friendId})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "friend_repository.Delete",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("связь между пользователем %s и другом %s не найдена", userId, friendId)
	}

	return nil
}

// GetFriends возвращает список ID друзей заданного пользователя.
func (r *repo) GetFriends(ctx context.Context, userId string) ([]string, error) {
	builder := sq.Select(friendColumn).
		From(tableName).
		Where(sq.Eq{idColumn: userId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "friend_repository.GetFriends",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []string
	for rows.Next() {
		var friendId string
		if err := rows.Scan(&friendId); err != nil {
			return nil, err
		}
		friends = append(friends, friendId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}
