package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetStressSamples(ctx context.Context, userID uuid.UUID) ([]entities.StressSampleExportEntity, error) {
	query := `
		SELECT id, timestamp_utc, window_minutes, encrypted_mean_hr, encrypted_rmssd_ms,
			encrypted_resting_hr, encrypted_steps, encrypted_sleep_debt_hours, created_at
		FROM stress_samples
		WHERE user_id = $1
		ORDER BY timestamp_utc ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.StressSampleExportEntity
	for rows.Next() {
		var ent entities.StressSampleExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.TimestampUTC, &ent.WindowMinutes,
			&ent.EncryptedMeanHR, &ent.EncryptedRMSSDms,
			&ent.EncryptedRestingHR, &ent.EncryptedSteps,
			&ent.EncryptedSleepDebtHours, &ent.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

func (r *repository) GetStressScores(ctx context.Context, userID uuid.UUID) ([]entities.StressScoreExportEntity, error) {
	query := `
		SELECT id, score, category, model_version, computed_at
		FROM stress_scores
		WHERE user_id = $1
		ORDER BY computed_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.StressScoreExportEntity
	for rows.Next() {
		var ent entities.StressScoreExportEntity
		if err := rows.Scan(&ent.ID, &ent.Score, &ent.Category, &ent.ModelVersion, &ent.ComputedAt); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

func (r *repository) GetMoodEntries(ctx context.Context, userID uuid.UUID) ([]entities.MoodEntryExportEntity, error) {
	query := `
		SELECT id, date, mood_score, encrypted_notes, created_at, updated_at
		FROM mood_entries
		WHERE user_id = $1
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.MoodEntryExportEntity
	for rows.Next() {
		var ent entities.MoodEntryExportEntity
		if err := rows.Scan(&ent.ID, &ent.Date, &ent.MoodScore, &ent.EncryptedNotes, &ent.CreatedAt, &ent.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
