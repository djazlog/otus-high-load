-- +goose Up
-- +goose StatementBegin
CREATE TABLE friends
(
    user_id          uuid NOT NULL references users(id) on delete cascade,
    friend_id          uuid NOT NULL references users(id) on delete cascade,

    created_at   timestamp NOT NULL DEFAULT NOW(),
    updated_at   timestamp  DEFAULT NOW(),

    PRIMARY KEY  (user_id, friend_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friends;
-- +goose StatementEnd
