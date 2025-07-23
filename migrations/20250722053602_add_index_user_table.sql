-- +goose Up
-- +goose StatementBegin
CREATE INDEX if not exists idx_users_first_second_pattern_id ON users (first_name text_pattern_ops, second_name text_pattern_ops, id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists idx_users_first_second_pattern_id;
-- +goose StatementEnd
