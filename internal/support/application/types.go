package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
)

type Service interface {
	AddMember(ctx context.Context, veteranID domain.VeteranId, memberId domain.MemberId) (domain.Member, error)
	GetAll(ctx context.Context, id domain.VeteranId) ([]domain.Member, error)
}

type service struct {
	repo infrastructure.Repository
}

func NewService(repo infrastructure.Repository) Service {
	return &service{repo: repo}
}
