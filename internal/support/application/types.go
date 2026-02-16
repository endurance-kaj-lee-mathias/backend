package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
)

type Service interface {
	// AddMember adds a supporter to a veteran. Both veteran and supporter roles
	// are loaded from the database so that the DB remains the source of truth.
	AddMember(ctx context.Context, veteranID domain.VeteranId, memberId domain.MemberId) (domain.Member, error)
	GetAll(ctx context.Context, id domain.VeteranId) ([]domain.Member, error)
	GetAllByMember(ctx context.Context, id domain.MemberId) ([]domain.Member, error)
	DeleteSupporter(ctx context.Context, veteranID domain.VeteranId, supportID domain.MemberId) error
}

type service struct {
	repo         infrastructure.Repository
	userRoleRead infrastructure.UserRoleReader
}

func NewService(repo infrastructure.Repository, userRoleRead infrastructure.UserRoleReader) Service {
	return &service{repo: repo, userRoleRead: userRoleRead}
}
