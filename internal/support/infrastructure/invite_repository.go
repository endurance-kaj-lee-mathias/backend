package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
)

const inviteQuery = `
	SELECT
		i.id, i.sender_id,
		s.encrypted_username, s.encrypted_first_name, s.encrypted_last_name, s.encrypted_user_key, s.image,
		i.receiver_id,
		r.encrypted_username, r.encrypted_first_name, r.encrypted_last_name, r.encrypted_user_key, r.image,
		i.status, i.note, i.created_at, i.updated_at
	FROM support_invites i
	JOIN users s ON s.id = i.sender_id
	JOIN users r ON r.id = i.receiver_id
`

func (r *inviteRepository) scanInvite(row *sql.Row) (entities.InviteEntity, error) {
	var ent entities.InviteEntity
	err := row.Scan(
		&ent.ID, &ent.SenderID,
		&ent.SenderEncryptedUsername, &ent.SenderEncryptedFirst, &ent.SenderEncryptedLast, &ent.SenderEncryptedUserKey, &ent.SenderImage,
		&ent.ReceiverID,
		&ent.ReceiverEncryptedUsername, &ent.ReceiverEncryptedFirst, &ent.ReceiverEncryptedLast, &ent.ReceiverEncryptedUserKey, &ent.ReceiverImage,
		&ent.Status, &ent.Note, &ent.CreatedAt, &ent.UpdatedAt,
	)
	return ent, err
}

func (r *inviteRepository) CreateInvite(ctx context.Context, inv domain.Invite) error {
	query := `
		INSERT INTO support_invites (id, sender_id, receiver_id, status, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		inv.ID.UUID, inv.Sender.ID.UUID, inv.Receiver.ID.UUID,
		string(inv.Status), inv.Note, inv.CreatedAt, inv.UpdatedAt,
	)
	return err
}

func (r *inviteRepository) FindInviteByID(ctx context.Context, id uuid.UUID) (entities.InviteEntity, error) {
	row := r.db.QueryRowContext(ctx, inviteQuery+` WHERE i.id = $1`, id)
	ent, err := r.scanInvite(row)
	if errors.Is(err, sql.ErrNoRows) {
		return entities.InviteEntity{}, InviteNotFound
	}
	return ent, err
}

func (r *inviteRepository) FindPendingBySenderReceiver(ctx context.Context, senderID, receiverID uuid.UUID) (entities.InviteEntity, bool, error) {
	row := r.db.QueryRowContext(ctx, inviteQuery+` WHERE i.sender_id = $1 AND i.receiver_id = $2 AND i.status = 'PENDING'`, senderID, receiverID)
	ent, err := r.scanInvite(row)
	if errors.Is(err, sql.ErrNoRows) {
		return entities.InviteEntity{}, false, nil
	}
	if err != nil {
		return entities.InviteEntity{}, false, err
	}
	return ent, true, nil
}

func (r *inviteRepository) DeleteInvite(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM support_invites WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return InviteNotFound
	}
	return nil
}

func (r *inviteRepository) ListPendingForUser(ctx context.Context, userID uuid.UUID) ([]entities.InviteEntity, error) {
	query := inviteQuery + ` WHERE (i.sender_id = $1 OR i.receiver_id = $1) AND i.status = 'PENDING'`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var ents []entities.InviteEntity
	for rows.Next() {
		var ent entities.InviteEntity
		if err := rows.Scan(
			&ent.ID, &ent.SenderID,
			&ent.SenderEncryptedUsername, &ent.SenderEncryptedFirst, &ent.SenderEncryptedLast, &ent.SenderEncryptedUserKey, &ent.SenderImage,
			&ent.ReceiverID,
			&ent.ReceiverEncryptedUsername, &ent.ReceiverEncryptedFirst, &ent.ReceiverEncryptedLast, &ent.ReceiverEncryptedUserKey, &ent.ReceiverImage,
			&ent.Status, &ent.Note, &ent.CreatedAt, &ent.UpdatedAt,
		); err != nil {
			return nil, err
		}
		ents = append(ents, ent)
	}

	return ents, rows.Err()
}
