package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetSentInvites(ctx context.Context, userID uuid.UUID) ([]entities.InviteExportEntity, error) {
	query := `
		SELECT i.id, i.receiver_id,
			r.encrypted_username, r.encrypted_first_name, r.encrypted_last_name, r.encrypted_user_key, r.image,
			i.status, i.note, i.created_at, i.updated_at
		FROM support_invites i
		JOIN users r ON r.id = i.receiver_id
		WHERE i.sender_id = $1
		ORDER BY i.created_at DESC
	`

	return r.queryInvites(ctx, query, userID)
}

func (r *repository) GetReceivedInvites(ctx context.Context, userID uuid.UUID) ([]entities.InviteExportEntity, error) {
	query := `
		SELECT i.id, i.sender_id,
			s.encrypted_username, s.encrypted_first_name, s.encrypted_last_name, s.encrypted_user_key, s.image,
			i.status, i.note, i.created_at, i.updated_at
		FROM support_invites i
		JOIN users s ON s.id = i.sender_id
		WHERE i.receiver_id = $1
		ORDER BY i.created_at DESC
	`

	return r.queryInvites(ctx, query, userID)
}

func (r *repository) queryInvites(ctx context.Context, query string, userID uuid.UUID) ([]entities.InviteExportEntity, error) {
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.InviteExportEntity
	for rows.Next() {
		var ent entities.InviteExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.OtherUserID,
			&ent.OtherEncryptedUsername, &ent.OtherEncryptedFirstName, &ent.OtherEncryptedLastName,
			&ent.OtherEncryptedUserKey, &ent.OtherImage,
			&ent.Status, &ent.Note, &ent.CreatedAt, &ent.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
