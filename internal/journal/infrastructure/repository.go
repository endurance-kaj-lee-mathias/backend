package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/infrastructure/entities"
)

func (r *repository) GetUserProfile(ctx context.Context, userID uuid.UUID) (entities.UserProfileEntity, error) {
	query := `
		SELECT id, encrypted_user_key, encrypted_first_name, encrypted_last_name, encrypted_username,
		       encrypted_about, encrypted_introduction, encrypted_phone_number, encrypted_roles, image, is_private
		FROM users
		WHERE id = $1
	`

	var ent entities.UserProfileEntity
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&ent.ID,
		&ent.EncryptedUserKey,
		&ent.EncryptedFirstName,
		&ent.EncryptedLastName,
		&ent.EncryptedUsername,
		&ent.EncryptedAbout,
		&ent.EncryptedIntroduction,
		&ent.EncryptedPhoneNumber,
		&ent.EncryptedRoles,
		&ent.Image,
		&ent.IsPrivate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.UserProfileEntity{}, UserNotFound
		}
		return entities.UserProfileEntity{}, err
	}

	return ent, nil
}

func (r *repository) GetStressScoresPaginated(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entities.StressScoreRow, int, error) {
	query := `
		SELECT id, score, category, model_version, computed_at, COUNT(*) OVER() AS total
		FROM stress_scores
		WHERE user_id = $1
		ORDER BY computed_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("journal: close stress scores rows", "error", err)
		}
	}()

	var result []entities.StressScoreRow
	var total int

	for rows.Next() {
		var row entities.StressScoreRow
		if err := rows.Scan(&row.ID, &row.Score, &row.Category, &row.ModelVersion, &row.ComputedAt, &total); err != nil {
			return nil, 0, err
		}
		result = append(result, row)
	}

	return result, total, rows.Err()
}

func (r *repository) GetMoodEntriesPaginated(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entities.MoodEntryRow, int, error) {
	query := `
		SELECT id, date, mood_score, encrypted_notes, created_at, updated_at, COUNT(*) OVER() AS total
		FROM mood_entries
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("journal: close mood entries rows", "error", err)
		}
	}()

	var result []entities.MoodEntryRow
	var total int

	for rows.Next() {
		var row entities.MoodEntryRow
		if err := rows.Scan(&row.ID, &row.Date, &row.MoodScore, &row.EncryptedNotes, &row.CreatedAt, &row.UpdatedAt, &total); err != nil {
			return nil, 0, err
		}
		result = append(result, row)
	}

	return result, total, rows.Err()
}
