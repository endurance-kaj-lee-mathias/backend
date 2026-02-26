package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/infrastructure/entities"
)

func (r *repository) FindConversation(ctx context.Context, userA, userB uuid.UUID) (entities.ConversationEntity, error) {
	query := `
		SELECT c.id, c.created_at
		FROM conversations c
		JOIN conversation_participants pa ON pa.conversation_id = c.id AND pa.user_id = $1
		JOIN conversation_participants pb ON pb.conversation_id = c.id AND pb.user_id = $2
		LIMIT 1
	`

	var ent entities.ConversationEntity
	err := r.db.QueryRowContext(ctx, query, userA, userB).Scan(&ent.ID, &ent.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ConversationEntity{}, ConversationNotFound
		}
		return entities.ConversationEntity{}, err
	}

	return ent, nil
}

func (r *repository) CreateConversation(ctx context.Context, ent entities.ConversationEntity) error {
	query := `INSERT INTO conversations (id, created_at) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, ent.ID, ent.CreatedAt)
	return err
}

func (r *repository) SaveParticipantKey(ctx context.Context, ent entities.ParticipantKeyEntity) error {
	query := `
		INSERT INTO conversation_participants (conversation_id, user_id, encrypted_conversation_key)
		VALUES ($1, $2, $3)
		ON CONFLICT (conversation_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, ent.ConversationID, ent.UserID, ent.EncryptedConversationKey)
	return err
}

func (r *repository) UpdateParticipantKey(ctx context.Context, conversationID, userID uuid.UUID, encryptedKey []byte) error {
	query := `
		UPDATE conversation_participants
		SET encrypted_conversation_key = $3
		WHERE conversation_id = $1 AND user_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, conversationID, userID, encryptedKey)
	return err
}

func (r *repository) GetUserEncryptedKey(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	query := `SELECT encrypted_user_key FROM users WHERE id = $1`

	var encryptedUserKey []byte
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&encryptedUserKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}

	return encryptedUserKey, nil
}

func (r *repository) GetParticipantKey(ctx context.Context, conversationID, userID uuid.UUID) (entities.ParticipantKeyEntity, error) {
	query := `
		SELECT cp.conversation_id, cp.user_id, cp.encrypted_conversation_key, u.encrypted_user_key
		FROM conversation_participants cp
		JOIN users u ON u.id = cp.user_id
		WHERE cp.conversation_id = $1 AND cp.user_id = $2
	`

	var ent entities.ParticipantKeyEntity
	err := r.db.QueryRowContext(ctx, query, conversationID, userID).Scan(
		&ent.ConversationID,
		&ent.UserID,
		&ent.EncryptedConversationKey,
		&ent.EncryptedUserKey,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ParticipantKeyEntity{}, ParticipantNotFound
		}
		return entities.ParticipantKeyEntity{}, err
	}

	return ent, nil
}

func (r *repository) CreateMessage(ctx context.Context, ent entities.MessageEntity) error {
	query := `
		INSERT INTO messages (id, conversation_id, sender_id, encrypted_content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, ent.ID, ent.ConversationID, ent.SenderID, ent.EncryptedContent, ent.CreatedAt)
	return err
}

func (r *repository) GetMessages(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]entities.MessageEntity, error) {
	query := `
		SELECT id, conversation_id, sender_id, encrypted_content, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var ents []entities.MessageEntity

	for rows.Next() {
		var ent entities.MessageEntity
		if err := rows.Scan(&ent.ID, &ent.ConversationID, &ent.SenderID, &ent.EncryptedContent, &ent.CreatedAt); err != nil {
			return nil, err
		}
		ents = append(ents, ent)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ents, nil
}

func (r *repository) CheckSupportRelationship(ctx context.Context, userA, userB uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_supports
			WHERE (veteran_id = $1 AND support_id = $2)
			   OR (veteran_id = $2 AND support_id = $1)
		)
	`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, userA, userB).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
