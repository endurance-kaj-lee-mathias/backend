package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure/entities"
)

func (r *repository) Create(ctx context.Context, ent entities.StressSampleEntity) error {
	query := `
		INSERT INTO stress_samples (id, user_id, timestamp_utc, window_minutes, encrypted_mean_hr, encrypted_rmssd_ms, encrypted_resting_hr, encrypted_steps, encrypted_sleep_debt_hours, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ent.ID,
		ent.UserID,
		ent.TimestampUTC,
		ent.WindowMinutes,
		ent.EncryptedMeanHR,
		ent.EncryptedRMSSDms,
		ent.EncryptedRestingHR,
		ent.EncryptedSteps,
		ent.EncryptedSleepDebtHours,
		ent.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" && strings.Contains(pgErr.ConstraintName, "user_id") {
			return UserNotFound
		}
		return err
	}

	return nil
}

func (r *repository) CountSamples(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM stress_samples WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *repository) GetLatestSampleTimestamp(ctx context.Context, userID uuid.UUID) (time.Time, error) {
	var ts time.Time
	err := r.db.QueryRowContext(
		ctx,
		`SELECT timestamp_utc FROM stress_samples WHERE user_id = $1 ORDER BY timestamp_utc DESC LIMIT 1`,
		userID,
	).Scan(&ts)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, SampleNotFound
		}
		return time.Time{}, err
	}
	return ts, nil
}

func (r *repository) GetSamplesLast90Days(ctx context.Context, userID uuid.UUID) ([]entities.StressSampleEntity, error) {
	query := `
		SELECT id, user_id, timestamp_utc, window_minutes, encrypted_mean_hr, encrypted_rmssd_ms,
		       encrypted_resting_hr, encrypted_steps, encrypted_sleep_debt_hours, created_at
		FROM stress_samples
		WHERE user_id = $1
		  AND created_at >= $2
		ORDER BY timestamp_utc ASC
	`

	cutoff := time.Now().UTC().AddDate(0, 0, -90)
	rows, err := r.db.QueryContext(ctx, query, userID, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.StressSampleEntity
	for rows.Next() {
		var ent entities.StressSampleEntity
		if err := rows.Scan(
			&ent.ID,
			&ent.UserID,
			&ent.TimestampUTC,
			&ent.WindowMinutes,
			&ent.EncryptedMeanHR,
			&ent.EncryptedRMSSDms,
			&ent.EncryptedRestingHR,
			&ent.EncryptedSteps,
			&ent.EncryptedSleepDebtHours,
			&ent.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

func (r *repository) CreateScore(ctx context.Context, ent entities.StressScoreEntity) error {
	query := `
		INSERT INTO stress_scores (id, user_id, score, category, model_version, computed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		ent.ID,
		ent.UserID,
		ent.Score,
		ent.Category,
		ent.ModelVersion,
		ent.ComputedAt,
	)
	return err
}

func (r *repository) GetLatestScore(ctx context.Context, userID uuid.UUID) (domain.StressScore, error) {
	query := `
		SELECT id, user_id, score, category, model_version, computed_at
		FROM stress_scores
		WHERE user_id = $1
		ORDER BY computed_at DESC
		LIMIT 1
	`

	var ent entities.StressScoreEntity
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&ent.ID,
		&ent.UserID,
		&ent.Score,
		&ent.Category,
		&ent.ModelVersion,
		&ent.ComputedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.StressScore{}, ScoreNotFound
		}
		return domain.StressScore{}, err
	}

	return entities.ScoreFromEntity(ent), nil
}
