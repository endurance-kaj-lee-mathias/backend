-- +goose Up
ALTER TABLE availability_slots ADD COLUMN series_id UUID NULL;
CREATE INDEX idx_slots_series_id ON availability_slots (series_id) WHERE series_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_slots_series_id;
ALTER TABLE availability_slots DROP COLUMN IF EXISTS series_id;
