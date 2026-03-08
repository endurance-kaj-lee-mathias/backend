-- +goose Up
-- +goose StatementBegin
ALTER TABLE support_invites
    ADD COLUMN note VARCHAR(300);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE support_invites
    DROP COLUMN note;
-- +goose StatementEnd

