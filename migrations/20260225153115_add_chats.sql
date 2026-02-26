-- +goose Up
-- +goose StatementBegin
CREATE TABLE conversations
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE conversation_participants
(
    conversation_id           UUID  NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    user_id                   UUID  NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_conversation_key BYTEA NOT NULL,
    PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX idx_conv_participants_user_id ON conversation_participants (user_id);

CREATE TABLE messages
(
    id              UUID PRIMARY KEY,
    conversation_id UUID                     NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    sender_id       UUID                     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_content BYTEA                  NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_messages_conversation_id ON messages (conversation_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversation_participants;
DROP TABLE IF EXISTS conversations;
-- +goose StatementEnd
