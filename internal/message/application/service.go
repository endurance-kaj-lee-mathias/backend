package application

import (
	"context"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/infrastructure/entities"
	"log/slog"
)

func (s *service) GetMessage(ctx context.Context) (domain.Message, error) {
	msg, err := s.repo.ReadMessage()

	if err != nil {
		slog.Error("message could not be fetched", "error", err)
		return domain.Message{}, err
	}

	return entities.FromEntity(msg)
}
