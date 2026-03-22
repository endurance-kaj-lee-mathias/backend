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
		       encrypted_about, encrypted_introduction, encrypted_phone_number, encrypted_roles, image, is_private, risk_level
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
		&ent.RiskLevel,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.UserProfileEntity{}, UserNotFound
		}
		return entities.UserProfileEntity{}, err
	}

	return ent, nil
}

func (r *repository) GetWeeklyAverages(ctx context.Context, userID uuid.UUID, weekOffset int) ([]entities.DailyAverageRow, error) {
	query := `
		SELECT
		    me.date,
		    AVG(me.mood_score)::float8                                          AS avg_mood,
		    AVG(ss.score)::float8                                               AS avg_stress,
		    (SELECT COUNT(DISTINCT date_trunc('week', date))
		     FROM mood_entries WHERE user_id = $1)                              AS total
		FROM mood_entries me
		LEFT JOIN stress_scores ss
		    ON ss.user_id = me.user_id
		    AND date_trunc('day', ss.computed_at AT TIME ZONE 'UTC') = me.date
		WHERE me.user_id = $1
		  AND date_trunc('week', me.date) = date_trunc('week', NOW() - ($2 * INTERVAL '1 week'))
		GROUP BY me.date
		ORDER BY me.date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, weekOffset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("journal: close weekly averages rows", "error", err)
		}
	}()

	var result []entities.DailyAverageRow

	for rows.Next() {
		var row entities.DailyAverageRow
		if err := rows.Scan(&row.Date, &row.AvgMood, &row.AvgStress, &row.Total); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, rows.Err()
}

func (r *repository) GetWeeklyMoodNotes(ctx context.Context, userID uuid.UUID, weekOffset int) ([]entities.MoodEntryNoteRow, error) {
	query := `
		SELECT date, encrypted_notes
		FROM mood_entries
		WHERE user_id = $1
		  AND date_trunc('week', date) = date_trunc('week', NOW() - ($2 * INTERVAL '1 week'))
		  AND encrypted_notes IS NOT NULL
		  AND octet_length(encrypted_notes) > 0
		ORDER BY date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, weekOffset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("journal: close weekly mood notes rows", "error", err)
		}
	}()

	var result []entities.MoodEntryNoteRow

	for rows.Next() {
		var row entities.MoodEntryNoteRow
		if err := rows.Scan(&row.Date, &row.EncryptedNotes); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, rows.Err()
}
