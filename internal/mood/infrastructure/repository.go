package infrastructure

import (
	"context"
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
		ON CONFLICT (user_id, date) DO UPDATE
		SET mood_score       = EXCLUDED.mood_score,
		    encrypted_notes  = EXCLUDED.encrypted_notes,
		    updated_at       = EXCLUDED.updated_at
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
