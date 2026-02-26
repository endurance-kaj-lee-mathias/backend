-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN role_hash TEXT NOT NULL DEFAULT '';

CREATE INDEX idx_users_role_hash ON users (role_hash);

CREATE TABLE mood_entries
(
    id               UUID                     NOT NULL PRIMARY KEY,
    user_id          UUID                     NOT NULL,
    date             DATE                     NOT NULL,
    mood_score       INT                      NOT NULL,
    encrypted_notes  BYTEA,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE (user_id, date),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_mood_entries_user_id ON mood_entries (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mood_entries;
DROP INDEX IF EXISTS idx_users_role_hash;
ALTER TABLE users DROP COLUMN IF EXISTS role_hash;
-- +goose StatementEnd

