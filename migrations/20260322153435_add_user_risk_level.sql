-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN risk_level VARCHAR(50) DEFAULT 'normal' NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN risk_level;
-- +goose StatementEnd
