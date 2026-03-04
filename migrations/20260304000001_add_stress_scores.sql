-- +goose Up
-- +goose StatementBegin
CREATE TABLE stress_scores
(
    id            UUID PRIMARY KEY,
    user_id       UUID                     NOT NULL,
    score         DOUBLE PRECISION         NOT NULL,
    category      VARCHAR(50)              NOT NULL,
    model_version VARCHAR(50)              NOT NULL,
    computed_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_stress_scores_user_id_computed_at ON stress_scores (user_id, computed_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stress_scores;
-- +goose StatementEnd

