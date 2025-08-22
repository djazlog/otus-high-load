-- +goose Up
-- +goose StatementBegin

-- Создание таблицы материализованных лент
CREATE TABLE IF NOT EXISTS materialized_feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    author_id UUID NOT NULL,
    post_text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, post_id)
);

-- Создание индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_materialized_feeds_user_id ON materialized_feeds(user_id);
CREATE INDEX IF NOT EXISTS idx_materialized_feeds_created_at ON materialized_feeds(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_materialized_feeds_author_id ON materialized_feeds(author_id);

-- Создание таблицы заданий на материализацию ленты
CREATE TABLE IF NOT EXISTS feed_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    error TEXT
);

-- Создание индексов для заданий
CREATE INDEX IF NOT EXISTS idx_feed_jobs_status ON feed_jobs(status);
CREATE INDEX IF NOT EXISTS idx_feed_jobs_priority ON feed_jobs(priority);
CREATE INDEX IF NOT EXISTS idx_feed_jobs_created_at ON feed_jobs(created_at);
CREATE INDEX IF NOT EXISTS idx_feed_jobs_user_id ON feed_jobs(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS materialized_feeds;
DROP TABLE IF EXISTS feed_jobs; 
-- +goose StatementEnd
