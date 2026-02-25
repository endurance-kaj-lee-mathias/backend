-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY,
    email      TEXT          NOT NULL UNIQUE,
    roles      JSONB         NOT NULL,
    created_at TIMESTAMP
                   WITH
                   TIME ZONE NOT NULL,
    updated_at TIMESTAMP
                   WITH
                   TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_supports
(
    veteran_id UUID          NOT NULL,
    support_id UUID          NOT NULL,
    created_at TIMESTAMP
                   WITH
                   TIME ZONE NOT NULL,
    PRIMARY KEY (veteran_id, support_id),
    FOREIGN KEY (veteran_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (support_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

DROP TABLE user_supports;

-- +goose StatementEnd
