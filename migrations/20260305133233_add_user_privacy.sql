-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN is_private BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS is_private;
-- +goose StatementEnd
