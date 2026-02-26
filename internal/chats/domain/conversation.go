package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Conversation struct {
	ID           ConversationId
	Participants []uuid.UUID
	CreatedAt    time.Time
}

func NewConversation(id ConversationId, participants []uuid.UUID, createdAt time.Time) Conversation {
	return Conversation{
		ID:           id,
		Participants: participants,
		CreatedAt:    createdAt,
	}
}
