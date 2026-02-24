package infrastructure

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
)

func (r *repository) Create(ctx context.Context, veteranID, memberId uuid.UUID) (entities.MemberEntity, error) {
	var ent entities.MemberEntity

	query := `
		WITH ins AS (
			INSERT INTO user_supports (veteran_id, support_id, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (veteran_id, support_id) DO NOTHING
			RETURNING veteran_id, support_id, created_at
		)
		SELECT ins.veteran_id, ins.support_id, ins.created_at,
			u.encrypted_email, u.encrypted_first_name, u.encrypted_last_name,
			u.encrypted_user_key, u.updated_at
		FROM users u
		LEFT JOIN ins ON u.id = ins.support_id
		LEFT JOIN user_supports s ON u.id = s.support_id AND s.veteran_id = $1
		WHERE u.id = $2
	`

	err := r.db.QueryRowContext(ctx, query, veteranID, memberId, time.Now().UTC()).Scan(
		&ent.Veteran, &ent.ID, &ent.CreatedAt,
		&ent.EncryptedEmail, &ent.EncryptedFirst, &ent.EncryptedLast,
		&ent.EncryptedUserKey, &ent.UpdatedAt,
	)

	if err != nil {
		return entities.MemberEntity{}, err
	}

	return ent, nil
}

func (r *repository) ReadAll(ctx context.Context, id uuid.UUID) ([]entities.MemberEntity, error) {
	query := `
		SELECT u.id, s.veteran_id,
			u.encrypted_email, u.encrypted_first_name, u.encrypted_last_name,
			u.encrypted_user_key, s.created_at, u.updated_at
		FROM users u
		JOIN user_supports s ON u.id = s.support_id
		WHERE s.veteran_id = $1
		UNION
		SELECT u.id, s.veteran_id,
			u.encrypted_email, u.encrypted_first_name, u.encrypted_last_name,
			u.encrypted_user_key, s.created_at, u.updated_at
		FROM users u
		JOIN user_supports s ON u.id = s.veteran_id
		WHERE s.support_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var ents []entities.MemberEntity

	for rows.Next() {
		var ent entities.MemberEntity

		err := rows.Scan(
			&ent.ID, &ent.Veteran,
			&ent.EncryptedEmail, &ent.EncryptedFirst, &ent.EncryptedLast,
			&ent.EncryptedUserKey, &ent.CreatedAt, &ent.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		ents = append(ents, ent)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ents, nil
}

func (r *repository) ReadAllByMember(ctx context.Context, id uuid.UUID) ([]entities.MemberEntity, error) {
	query := `
		SELECT u.id, s.veteran_id,
			u.encrypted_email, u.encrypted_first_name, u.encrypted_last_name,
			u.encrypted_user_key, s.created_at, u.updated_at
		FROM users u
		JOIN user_supports s ON u.id = s.veteran_id
		WHERE s.support_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var ents []entities.MemberEntity

	for rows.Next() {
		var ent entities.MemberEntity

		err := rows.Scan(
			&ent.ID, &ent.Veteran,
			&ent.EncryptedEmail, &ent.EncryptedFirst, &ent.EncryptedLast,
			&ent.EncryptedUserKey, &ent.CreatedAt, &ent.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		ents = append(ents, ent)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ents, nil
}

func (r *repository) Delete(ctx context.Context, veteranID, supportID uuid.UUID) error {
	query := `DELETE FROM user_supports WHERE veteran_id = $1 AND support_id = $2`

	result, err := r.db.ExecContext(ctx, query, veteranID, supportID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return UserNotFound
	}

	return nil
}
