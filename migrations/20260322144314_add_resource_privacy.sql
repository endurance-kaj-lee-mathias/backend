-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resource_privacy
(
    owner_id   UUID    NOT NULL,
    resource   TEXT    NOT NULL,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,

    PRIMARY KEY (owner_id, resource),
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resource_privacy;
-- +goose StatementEnd
