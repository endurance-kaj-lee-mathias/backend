-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS user_supports;
DROP TABLE IF EXISTS user_addresses;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id                     UUID PRIMARY KEY,
    email_hash             TEXT NOT NULL UNIQUE,
    phone_number_hash      TEXT,
    encrypted_email        BYTEA NOT NULL,
    encrypted_first_name   BYTEA NOT NULL,
    encrypted_last_name    BYTEA NOT NULL,
    encrypted_phone_number BYTEA,
    encrypted_roles        BYTEA NOT NULL,
    encrypted_user_key     BYTEA NOT NULL,
    key_version            INT NOT NULL DEFAULT 1,
    created_at             TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at             TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE user_addresses (
    id                       UUID PRIMARY KEY,
    user_id                  UUID NOT NULL UNIQUE,
    encrypted_street         BYTEA NOT NULL,
    encrypted_house_number   BYTEA NOT NULL,
    encrypted_postal_code    BYTEA NOT NULL,
    encrypted_city           BYTEA NOT NULL,
    encrypted_country        BYTEA NOT NULL,
    created_at               TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE user_supports (
    veteran_id UUID NOT NULL,
    support_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (veteran_id, support_id),
    FOREIGN KEY (veteran_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (support_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_supports;
DROP TABLE IF EXISTS user_addresses;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id          UUID PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    first_name  TEXT NOT NULL DEFAULT '',
    last_name   TEXT NOT NULL DEFAULT '',
    phone_number TEXT,
    roles       JSONB NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE user_addresses (
    id           UUID PRIMARY KEY,
    user_id      UUID NOT NULL UNIQUE,
    street       TEXT NOT NULL,
    house_number TEXT NOT NULL,
    postal_code  TEXT NOT NULL,
    city         TEXT NOT NULL,
    country      TEXT NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE user_supports (
    veteran_id UUID NOT NULL,
    support_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (veteran_id, support_id),
    FOREIGN KEY (veteran_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (support_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

