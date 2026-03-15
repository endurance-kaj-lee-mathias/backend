package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
)

type MessageEntity struct {
	ID               uuid.UUID `db:"id"`
	ConversationID   uuid.UUID `db:"conversation_id"`
	SenderID         uuid.UUID `db:"sender_id"`
	SenderUsername   []byte    `db:"encrypted_username"`
	SenderUserKey    []byte    `db:"encrypted_user_key"`
	EncryptedContent []byte    `db:"encrypted_content"`
	CreatedAt        time.Time `db:"created_at"`
}

func FromMessageEntity(ent MessageEntity, convKey []byte, enc encryption.Service) (domain.Message, error) {
	senderKey, err := enc.DecryptUserKey(ent.SenderUserKey)
	if err != nil {
		return domain.Message{}, err
	}

	usernameBytes, err := enc.Decrypt(ent.SenderUsername, senderKey)
	if err != nil {
		return domain.Message{}, err
	}

	contentBytes, err := enc.Decrypt(ent.EncryptedContent, convKey)
	if err != nil {
		return domain.Message{}, err
	}

	convID := domain.ConversationId{UUID: ent.ConversationID}
	msgID := domain.MessageId{UUID: ent.ID}

	return domain.NewMessage(msgID, convID, ent.SenderID, string(usernameBytes), string(contentBytes), ent.CreatedAt), nil
}

func FromMessageEntities(ents []MessageEntity, convKey []byte, enc encryption.Service) ([]domain.Message, error) {
	out := make([]domain.Message, 0, len(ents))
	for _, ent := range ents {
		msg, err := FromMessageEntity(ent, convKey, enc)
		if err != nil {
			return nil, err
		}
		out = append(out, msg)
	}
	return out, nil
}
