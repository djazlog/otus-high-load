-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts
(
    id             uuid      NOT NULL,
    content        text,
    author_user_id uuid      NOT NULL references users (id) ON DELETE CASCADE,
    created_at     timestamp NOT NULL,
    updated_at     timestamp,

    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd
