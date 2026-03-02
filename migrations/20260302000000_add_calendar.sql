-- +goose Up
-- +goose StatementBegin
CREATE TABLE availability_slots
(
    id          UUID PRIMARY KEY,
    provider_id UUID                     NOT NULL,
    start_time  TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time    TIMESTAMP WITH TIME ZONE NOT NULL,
    is_urgent   BOOLEAN                  NOT NULL DEFAULT FALSE,
    is_booked   BOOLEAN                  NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (provider_id) REFERENCES users (id) ON DELETE CASCADE,
    CHECK (end_time > start_time)
);

CREATE UNIQUE INDEX idx_slots_no_overlap ON availability_slots (provider_id, start_time, end_time);
CREATE INDEX idx_slots_provider_start ON availability_slots (provider_id, start_time);
CREATE INDEX idx_slots_start_time ON availability_slots (start_time);
CREATE INDEX idx_slots_is_booked ON availability_slots (is_booked);

CREATE TABLE appointments
(
    id         UUID PRIMARY KEY,
    slot_id    UUID                     NOT NULL UNIQUE,
    veteran_id UUID                     NOT NULL,
    status     TEXT                     NOT NULL DEFAULT 'BOOKED',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (slot_id) REFERENCES availability_slots (id) ON DELETE CASCADE,
    FOREIGN KEY (veteran_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS availability_slots;
-- +goose StatementEnd

