-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS user_addresses;

CREATE TABLE user_addresses
(
    id                  UUID PRIMARY KEY,
    user_id             UUID                     NOT NULL UNIQUE,
    encrypted_street    BYTEA                    NOT NULL,
    encrypted_locality  BYTEA                    NOT NULL,
    encrypted_region    BYTEA                    NOT NULL,
    encrypted_postal_code BYTEA                  NOT NULL,
    encrypted_country   BYTEA                    NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_addresses;

CREATE TABLE user_addresses
(
    id                     UUID PRIMARY KEY,
    user_id                UUID                     NOT NULL UNIQUE,
    encrypted_street       BYTEA                    NOT NULL,
    encrypted_house_number BYTEA                    NOT NULL,
    encrypted_postal_code  BYTEA                    NOT NULL,
    encrypted_city         BYTEA                    NOT NULL,
    encrypted_country      BYTEA                    NOT NULL,
    created_at             TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd
