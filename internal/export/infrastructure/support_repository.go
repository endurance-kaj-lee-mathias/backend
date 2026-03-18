package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetSupporters(ctx context.Context, veteranID uuid.UUID) ([]entities.SupportMemberExportEntity, error) {
	query := `
		SELECT u.id, u.encrypted_email, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.image, u.encrypted_user_key, s.created_at
		FROM users u
		JOIN user_supports s ON u.id = s.support_id
		WHERE s.veteran_id = $1
		UNION
		SELECT u.id, u.encrypted_email, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.image, u.encrypted_user_key, s.created_at
		FROM users u
		JOIN user_supports s ON u.id = s.veteran_id
		WHERE s.support_id = $1
	`

	return r.querySupportMembers(ctx, query, veteranID)
}

func (r *repository) GetSupportedVeterans(ctx context.Context, supportID uuid.UUID) ([]entities.SupportMemberExportEntity, error) {
	query := `
		SELECT u.id, u.encrypted_email, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.image, u.encrypted_user_key, s.created_at
		FROM users u
		JOIN user_supports s ON u.id = s.veteran_id
		WHERE s.support_id = $1
		UNION
		SELECT u.id, u.encrypted_email, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.image, u.encrypted_user_key, s.created_at
		FROM users u
		JOIN user_supports s ON u.id = s.support_id
		WHERE s.veteran_id = $1
	`

	return r.querySupportMembers(ctx, query, supportID)
}

func (r *repository) querySupportMembers(ctx context.Context, query string, id uuid.UUID) ([]entities.SupportMemberExportEntity, error) {
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.SupportMemberExportEntity
	for rows.Next() {
		var ent entities.SupportMemberExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.EncryptedEmail, &ent.EncryptedUsername,
			&ent.EncryptedFirst, &ent.EncryptedLast,
			&ent.Image, &ent.EncryptedUserKey, &ent.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
