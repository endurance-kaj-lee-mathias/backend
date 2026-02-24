-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN username_hash       TEXT NOT NULL,
    ADD COLUMN encrypted_username  BYTEA NOT NULL;

CREATE UNIQUE INDEX idx_users_username_hash ON users (username_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_username_hash;

ALTER TABLE users
    DROP COLUMN IF EXISTS username_hash,
    DROP COLUMN IF EXISTS encrypted_username;
-- +goose StatementEnd

