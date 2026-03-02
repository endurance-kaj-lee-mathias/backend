-- +goose Up
-- +goose StatementBegin
CREATE TABLE support_invites
(
    id          UUID                     NOT NULL PRIMARY KEY,
    sender_id   UUID                     NOT NULL,
    receiver_id UUID                     NOT NULL,
    status      TEXT                     NOT NULL DEFAULT 'PENDING',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_support_invites_pending ON support_invites (sender_id, receiver_id) WHERE status = 'PENDING';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_support_invites_pending;
DROP TABLE IF EXISTS support_invites;
-- +goose StatementEnd

