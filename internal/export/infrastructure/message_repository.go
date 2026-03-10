package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetConversationsAndMessages(ctx context.Context, userID uuid.UUID) ([]entities.MessageExportEntity, error) {
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.encrypted_content, m.created_at,
			cp_user.encrypted_conversation_key,
			cp_other.user_id,
			u_other.encrypted_username, u_other.encrypted_first_name, u_other.encrypted_last_name,
			u_other.encrypted_user_key
		FROM messages m
		JOIN conversation_participants cp_user ON cp_user.conversation_id = m.conversation_id AND cp_user.user_id = $1
		JOIN conversation_participants cp_other ON cp_other.conversation_id = m.conversation_id AND cp_other.user_id != $1
		JOIN users u_other ON u_other.id = cp_other.user_id
		WHERE m.conversation_id IN (
			SELECT conversation_id FROM conversation_participants WHERE user_id = $1
		)
		ORDER BY m.created_at ASC
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

	var result []entities.MessageExportEntity
	for rows.Next() {
		var ent entities.MessageExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.ConversationID, &ent.SenderID, &ent.EncryptedContent, &ent.CreatedAt,
			&ent.EncryptedConversationKey,
			&ent.OtherParticipantID,
			&ent.OtherParticipantEncryptedUsername, &ent.OtherParticipantEncryptedFirstName,
			&ent.OtherParticipantEncryptedLastName, &ent.OtherParticipantEncryptedUserKey,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
