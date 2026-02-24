-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN encrypted_about        BYTEA,
    ADD COLUMN encrypted_introduction BYTEA,
    ADD COLUMN image                  TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN IF EXISTS encrypted_about,
    DROP COLUMN IF EXISTS encrypted_introduction,
    DROP COLUMN IF EXISTS image;
-- +goose StatementEnd

