package application

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"
)

type Service interface {
	CreateUser(ctx context.Context, email string, roles []domain.Role) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	AddSupportMember(ctx context.Context, veteranID, supportID uuid.UUID) error
	ListSupportMembers(ctx context.Context, veteranID uuid.UUID) ([]domain.User, error)
}

type service struct {
	repo infrastructure.Repository
}

func NewService(repo infrastructure.Repository) Service {
	return &service{repo: repo}
}
