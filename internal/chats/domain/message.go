package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID             MessageId
	ConversationID ConversationId
	SenderID       uuid.UUID
	Content        string
	CreatedAt      time.Time
}

func NewMessage(id MessageId, conversationID ConversationId, senderID uuid.UUID, content string, createdAt time.Time) Message {
	return Message{
		ID:             id,
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      createdAt,
	}
}
