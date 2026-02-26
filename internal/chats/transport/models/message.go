package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
)

type MessageModel struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversationId"`
	SenderID       uuid.UUID `json:"senderId"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

func ToMessageModel(m domain.Message) MessageModel {
	return MessageModel{
		ID:             m.ID.UUID,
		ConversationID: m.ConversationID.UUID,
		SenderID:       m.SenderID,
		Content:        m.Content,
		CreatedAt:      m.CreatedAt,
	}
}

func ToMessageModels(msgs []domain.Message) []MessageModel {
	out := make([]MessageModel, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, ToMessageModel(m))
	}
	return out
}
