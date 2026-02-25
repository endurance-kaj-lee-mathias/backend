package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure"
)

type Service interface {
	IngestSample(ctx context.Context, sample domain.StressSample) error
}

type service struct {
	repo          infrastructure.Repository
	userKeyReader infrastructure.UserKeyReader
	enc           encryption.Service
}

func NewService(repo infrastructure.Repository, userKeyReader infrastructure.UserKeyReader, enc encryption.Service) Service {
	return &service{repo: repo, userKeyReader: userKeyReader, enc: enc}
}
