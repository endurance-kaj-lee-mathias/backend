package infrastructure

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
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
