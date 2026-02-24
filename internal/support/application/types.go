package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
)

type Service interface {
	AddMember(ctx context.Context, veteranID domain.VeteranId, memberId domain.MemberId) (domain.Member, error)
	GetAll(ctx context.Context, id domain.VeteranId) ([]domain.Member, error)
	GetAllByMember(ctx context.Context, id domain.MemberId) ([]domain.Member, error)
	DeleteSupporter(ctx context.Context, veteranID domain.VeteranId, supportID domain.MemberId) error
}

type service struct {
	repo         infrastructure.Repository
	userRoleRead infrastructure.UserRoleReader
	enc          encryption.Service
}

func NewService(repo infrastructure.Repository, userRoleRead infrastructure.UserRoleReader, enc encryption.Service) Service {
	return &service{repo: repo, userRoleRead: userRoleRead, enc: enc}
}
