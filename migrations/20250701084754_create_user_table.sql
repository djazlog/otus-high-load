-- +goose Up
CREATE TABLE users
(
    id          uuid NOT NULL,
    first_name  text,
    second_name text,
    biography   text,
    city        text,
    password    text,
    birth_date   timestamp,
    created_at   timestamp NOT NULL,
    updated_at   timestamp,

    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;
