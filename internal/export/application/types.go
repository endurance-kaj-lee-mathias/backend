package application

import (
	"context"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure"
)

type Service interface {
	ExportUserData(ctx context.Context, userID uuid.UUID) (domain.UserDataExport, error)
}

type service struct {
	repo infrastructure.Repository
	enc  encryption.Service
}

func NewService(repo infrastructure.Repository, enc encryption.Service) Service {
	return &service{repo: repo, enc: enc}
}
