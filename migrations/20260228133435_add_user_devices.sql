-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_devices
(
    id           UUID                     NOT NULL PRIMARY KEY,
    user_id      UUID                     NOT NULL,
    device_token TEXT                     NOT NULL UNIQUE,
    platform     TEXT                     NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_devices_user_id ON user_devices (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_devices;
-- +goose StatementEnd

