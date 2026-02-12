package infrastructure

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/infrastructure/entities"
)

func (r *repository) ReadMessage() (entities.MessageEntity, error) {
	var msg entities.MessageEntity

	err := r.db.QueryRow("SELECT NOW()").Scan(
		&msg.Value,
	)

	if err != nil {
		return entities.MessageEntity{}, err
	}

	return msg, nil
}
