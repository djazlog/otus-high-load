-- +goose Up
-- +goose NO TRANSACTION
CREATE TABLE IF NOT EXISTS dialog_messages (
    id uuid NOT NULL,
    from_user_id uuid NOT NULL,
    to_user_id uuid NOT NULL,
    text TEXT NOT NULL,
    dialog_key uuid NOT NULL, -- ключ диалога (шард-ключ) для распределения по шардам
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

-- создаём уникальный индекс, включающий dialog_key
--CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS dialog_messages_dialog_id_uidx ON dialog_messages (dialog_key, id);
ALTER TABLE dialog_messages
    ADD CONSTRAINT dialog_messages_pkey PRIMARY KEY (dialog_key, id);

-- Распределяем таблицу по ключу диалога
SELECT
    create_distributed_table('dialog_messages', 'dialog_key');

-- последние сообщения в диалоге
CREATE INDEX IF NOT EXISTS dialog_messages_dialog_created_idx ON dialog_messages (dialog_key, created_at DESC);

-- поиск по участникам (если нужно)
CREATE INDEX IF NOT EXISTS dialog_messages_from_to_idx ON dialog_messages (from_user_id, to_user_id);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dialog_messages;

-- +goose StatementEnd
