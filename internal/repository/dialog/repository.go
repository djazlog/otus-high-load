package dialog

import (
	"context"
	"otus-project/internal/client/db"
	"otus-project/internal/model"
	"otus-project/internal/repository"
	"otus-project/internal/repository/dialog/converter"
	repoModel "otus-project/internal/repository/dialog/model"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

const (
	tableName = "dialog_messages"

	idColumn         = "id"
	fromUserIdColumn = "from_user_id"
	toUserIdColumn   = "to_user_id"
	textColumn       = "text"
	createdAtColumn  = "created_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.DialogRepository {
	return &repo{db: db}
}

// SendMessage сохраняет сообщение в диалоге
func (r *repo) SendMessage(ctx context.Context, fromUserId, toUserId, text string) error {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(fromUserIdColumn, toUserIdColumn, textColumn, createdAtColumn).
		Values(fromUserId, toUserId, text, time.Now())

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}

	q := db.Query{
		Name:     "dialog_repository.SendMessage",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert query")
	}

	return nil
}

// GetDialogList возвращает список сообщений диалога между двумя пользователями
func (r *repo) GetDialogList(ctx context.Context, userId1, userId2 string) ([]*model.DialogMessage, error) {
	// Получаем сообщения в обоих направлениях (от userId1 к userId2 и наоборот)
	builder := sq.Select(idColumn, fromUserIdColumn, toUserIdColumn, textColumn, createdAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(
			sq.Or{
				sq.And{
					sq.Eq{fromUserIdColumn: userId1},
					sq.Eq{toUserIdColumn: userId2},
				},
				sq.And{
					sq.Eq{fromUserIdColumn: userId2},
					sq.Eq{toUserIdColumn: userId1},
				},
			},
		).
		OrderBy(createdAtColumn + " ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}

	q := db.Query{
		Name:     "dialog_repository.GetDialogList",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute select query")
	}
	defer rows.Close()

	var messages []*repoModel.DialogMessage
	for rows.Next() {
		var msg repoModel.DialogMessage
		err := rows.Scan(&msg.ID, &msg.FromUserID, &msg.ToUserID, &msg.Text, &msg.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		messages = append(messages, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating rows")
	}

	// Конвертируем в сервисные модели
	return converter.ToDialogMessagesFromRepo(messages), nil
}
