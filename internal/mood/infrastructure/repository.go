package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure/entities"
)

func (r *repository) Upsert(ctx context.Context, ent entities.MoodEntryEntity) error {
	query := `
		INSERT INTO mood_entries (id, user_id, date, mood_score, encrypted_notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ent.ID,
		ent.UserID,
		ent.Date,
		ent.MoodScore,
		ent.EncryptedNotes,
		ent.CreatedAt,
		ent.UpdatedAt,
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

func (r *repository) FindVeteransWithoutEntryInLast24Hours(ctx context.Context, veteranRoleHash string) ([]uuid.UUID, error) {
	query := `
		SELECT u.id
		FROM users u
		LEFT JOIN mood_entries m
		    ON m.user_id = u.id
		    AND m.updated_at > NOW() - INTERVAL '24 hours'
		WHERE u.role_hash = $1
		  AND m.user_id IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, veteranRoleHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (r *repository) FindDeviceTokensByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `SELECT device_token FROM user_devices WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (r *repository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]entities.MoodEntryEntity, error) {
	query := `
		SELECT id, user_id, date, mood_score, encrypted_notes, created_at, updated_at
		FROM mood_entries
		WHERE user_id = $1
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.MoodEntryEntity

	for rows.Next() {
		var ent entities.MoodEntryEntity
		if err := rows.Scan(&ent.ID, &ent.UserID, &ent.Date, &ent.MoodScore, &ent.EncryptedNotes, &ent.CreatedAt, &ent.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

func (r *repository) FindTodayByUserID(ctx context.Context, userID uuid.UUID) (*entities.MoodEntryEntity, error) {
	query := `
		SELECT id, user_id, date, mood_score, encrypted_notes, created_at, updated_at
		FROM mood_entries
		WHERE user_id = $1
		  AND date = CURRENT_DATE
		LIMIT 1
	`

	var ent entities.MoodEntryEntity
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&ent.ID, &ent.UserID, &ent.Date, &ent.MoodScore, &ent.EncryptedNotes, &ent.CreatedAt, &ent.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, MoodEntryNotFound
		}
		return nil, err
	}

	return &ent, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*entities.MoodEntryEntity, error) {
	query := `
		SELECT id, user_id, date, mood_score, encrypted_notes, created_at, updated_at
		FROM mood_entries
		WHERE id = $1
	`

	var ent entities.MoodEntryEntity
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ent.ID, &ent.UserID, &ent.Date, &ent.MoodScore, &ent.EncryptedNotes, &ent.CreatedAt, &ent.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, MoodEntryNotFound
		}
		return nil, err
	}

	return &ent, nil
}

func (r *repository) Update(ctx context.Context, ent entities.MoodEntryEntity) error {
	query := `
		UPDATE mood_entries
		SET date = $1, mood_score = $2, encrypted_notes = $3, updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, ent.Date, ent.MoodScore, ent.EncryptedNotes, ent.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM mood_entries WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM mood_entries WHERE user_id = $1`, userID)
	return err
}
