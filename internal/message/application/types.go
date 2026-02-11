package application

import (
	"context"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/infrastructure"
)

type Service interface {
	GetMessage(
		ctx context.Context,
	) (domain.Message, error)
}

type service struct {
	repo infrastructure.Repository
}

func NewService(repo infrastructure.Repository) Service {
	return &service{repo}
}
