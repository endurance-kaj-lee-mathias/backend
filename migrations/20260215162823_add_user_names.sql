-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN first_name TEXT NOT NULL DEFAULT '';
ALTER TABLE users
    ADD COLUMN last_name TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN first_name;
ALTER TABLE users
    DROP COLUMN last_name;
-- +goose StatementEnd
