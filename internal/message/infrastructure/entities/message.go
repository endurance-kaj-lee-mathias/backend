package entities

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/domain"
)

type MessageEntity struct {
	Value string
}

func FromEntity(entity MessageEntity) (domain.Message, error) {
	return domain.Message{
		Value: entity.Value,
	}, nil
}
