-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS authorization_rules
(
    id         UUID PRIMARY KEY,
    owner_id   UUID          NOT NULL,
    viewer_id  UUID          NOT NULL,
    resource   TEXT          NOT NULL,
    effect     TEXT          NOT NULL,
    created_at TIMESTAMP
                   WITH
                   TIME ZONE NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (viewer_id) REFERENCES users (id) ON DELETE CASCADE,

    UNIQUE (owner_id, viewer_id, resource)
);

CREATE INDEX idx_authorization_rules_owner ON authorization_rules (owner_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS authorization_rules;
-- +goose StatementEnd
