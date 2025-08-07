-- +goose Up
-- +goose StatementBegin
CREATE TABLE friends
(
    user_id          uuid NOT NULL constraint fk_user_id references users(id) on delete cascade,
    friend_id          uuid NOT NULL constraint fk_friend_id references users(id) on delete cascade,

    created_at   timestamp NOT NULL DEFAULT NOW(),
    updated_at   timestamp  DEFAULT NOW(),

    UNIQUE (user_id, friend_id),

    PRIMARY KEY (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friends;
-- +goose StatementEnd
