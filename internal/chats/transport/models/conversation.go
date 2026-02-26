package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
)

type ConversationModel struct {
	ID           uuid.UUID   `json:"id"`
	Participants []uuid.UUID `json:"participants"`
	CreatedAt    time.Time   `json:"createdAt"`
}

func ToConversationModel(c domain.Conversation) ConversationModel {
	return ConversationModel{
		ID:           c.ID.UUID,
		Participants: c.Participants,
		CreatedAt:    c.CreatedAt,
	}
}
