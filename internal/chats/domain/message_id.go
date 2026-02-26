package domain

import (
	"github.com/gofrs/uuid"
)

type MessageId struct {
	uuid.UUID
}

func NewMessageId() MessageId {
	return MessageId{UUID: uuid.Must(uuid.NewV4())}
}

func ParseMessageId(value string) (MessageId, error) {
	id, err := uuid.FromString(value)
	if err != nil {
		return MessageId{}, err
	}
	return MessageId{UUID: id}, nil
}
