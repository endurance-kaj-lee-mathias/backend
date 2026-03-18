-- +goose Up
-- +goose StatementBegin
ALTER TABLE appointments ADD COLUMN title VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appointments DROP COLUMN title;
-- +goose StatementEnd
