package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"
)

type Service interface {
	AddUser(ctx context.Context, email string, roles []domain.Role) (domain.User, error)
	GetByID(ctx context.Context, id domain.UserId) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type service struct {
	repo infrastructure.Repository
}

func NewService(repo infrastructure.Repository) Service {
	return &service{repo: repo}
}
