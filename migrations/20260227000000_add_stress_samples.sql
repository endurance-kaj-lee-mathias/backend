-- +goose Up
-- +goose StatementBegin
CREATE TABLE stress_samples
(
    id                         UUID PRIMARY KEY,
    user_id                    UUID                     NOT NULL,
    timestamp_utc              TIMESTAMP WITH TIME ZONE NOT NULL,
    window_minutes             INT                      NOT NULL,
    encrypted_mean_hr          BYTEA                    NOT NULL,
    encrypted_rmssd_ms         BYTEA                    NOT NULL,
    encrypted_resting_hr       BYTEA,
    encrypted_steps            BYTEA,
    encrypted_sleep_debt_hours BYTEA,
    created_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_stress_samples_user_id ON stress_samples (user_id);
CREATE INDEX idx_stress_samples_timestamp ON stress_samples (user_id, timestamp_utc);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stress_samples;
-- +goose StatementEnd

