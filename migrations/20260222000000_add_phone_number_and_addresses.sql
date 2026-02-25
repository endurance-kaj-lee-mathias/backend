-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN phone_number TEXT;

CREATE TABLE IF NOT EXISTS user_addresses
(
    id           UUID PRIMARY KEY,
    user_id      UUID                     NOT NULL UNIQUE,
    street       TEXT                     NOT NULL,
    house_number TEXT                     NOT NULL,
    postal_code  TEXT                     NOT NULL,
    city         TEXT                     NOT NULL,
    country      TEXT                     NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_addresses;

ALTER TABLE users
    DROP COLUMN phone_number;
-- +goose StatementEnd

