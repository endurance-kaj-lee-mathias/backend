package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type MessageExportEntity struct {
	ID                                 uuid.UUID
	ConversationID                     uuid.UUID
	SenderID                           uuid.UUID
	EncryptedContent                   []byte
	CreatedAt                          time.Time
	EncryptedConversationKey           []byte
	OtherParticipantID                 uuid.UUID
	OtherParticipantEncryptedUsername  []byte
	OtherParticipantEncryptedFirstName []byte
	OtherParticipantEncryptedLastName  []byte
	OtherParticipantEncryptedUserKey   []byte
}
