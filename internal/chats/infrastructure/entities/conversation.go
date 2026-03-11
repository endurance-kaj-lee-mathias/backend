package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type ConversationEntity struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

type ConversationWithParticipantsEntity struct {
	ID           uuid.UUID `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	Participants []uuid.UUID
}

type ParticipantKeyEntity struct {
	ConversationID           uuid.UUID `db:"conversation_id"`
	UserID                   uuid.UUID `db:"user_id"`
	EncryptedConversationKey []byte    `db:"encrypted_conversation_key"`
	EncryptedUserKey         []byte    `db:"encrypted_user_key"`
}

type ConversationSummaryEntity struct {
	ConversationID                 uuid.UUID
	OtherUserID                    uuid.UUID
	OtherEncryptedFirstName        []byte
	OtherEncryptedLastName         []byte
	OtherEncryptedUserKey          []byte
	OtherImage                     *string
	CallerEncryptedConversationKey []byte
	CallerEncryptedUserKey         []byte
	LatestEncryptedContent         []byte
	LatestSenderID                 *uuid.UUID
	LatestMessageAt                *time.Time
}
