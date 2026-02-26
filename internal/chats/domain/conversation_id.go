package domain

import (
	"github.com/gofrs/uuid"
)

type ConversationId struct {
	uuid.UUID
}

func NewConversationId() ConversationId {
	return ConversationId{UUID: uuid.Must(uuid.NewV4())}
}

func ParseConversationId(value string) (ConversationId, error) {
	id, err := uuid.FromString(value)
	if err != nil {
		return ConversationId{}, err
	}
	return ConversationId{UUID: id}, nil
}
