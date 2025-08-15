-- +goose Up
-- +goose StatementBegin
-- Создание таблицы для сообщений диалогов
CREATE TABLE IF NOT EXISTS dialog_messages (
    id SERIAL PRIMARY KEY,
    from_user_id VARCHAR(255) NOT NULL,
    to_user_id VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX idx_dialog_messages_from_user_id ON dialog_messages (from_user_id);
CREATE INDEX idx_dialog_messages_to_user_id ON dialog_messages (to_user_id);
CREATE INDEX idx_dialog_messages_created_at ON dialog_messages (created_at);

-- Составной индекс для поиска диалога между двумя пользователями
CREATE INDEX idx_dialog_messages_users ON dialog_messages (from_user_id, to_user_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dialog_messages;
-- +goose StatementEnd

