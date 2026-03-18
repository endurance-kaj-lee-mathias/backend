package application

import (
	"context"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/infrastructure"
)

type Service interface {
	GetJournal(ctx context.Context, viewerID uuid.UUID, veteranID uuid.UUID, weekOffset int) (domain.JournalReport, error)
}

type AuthorizationChecker interface {
	IsAllowed(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID, resource string) (bool, error)
}

type service struct {
	repo  infrastructure.Repository
	enc   encryption.Service
	authz AuthorizationChecker
}

func NewService(repo infrastructure.Repository, enc encryption.Service, authz AuthorizationChecker) Service {
	return &service{repo: repo, enc: enc, authz: authz}
}
