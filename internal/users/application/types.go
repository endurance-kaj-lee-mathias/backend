package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"
)

type Service interface {
	GetOrCreate(ctx context.Context, id domain.UserId, email string, username string, firstName string, lastName string, roles []domain.Role) (domain.User, error)
	GetByID(ctx context.Context, id domain.UserId) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	DeleteUser(ctx context.Context, id domain.UserId) error
	UpdatePhoneNumber(ctx context.Context, id domain.UserId, phoneNumber *string) error
	UpdateIntroduction(ctx context.Context, id domain.UserId, introduction string) error
	UpdateAbout(ctx context.Context, id domain.UserId, about string) error
	UpdateImage(ctx context.Context, id domain.UserId, image string) error
	UpsertAddress(ctx context.Context, userID domain.UserId, street string, houseNumber string, postalCode string, city string, country string) (domain.Address, error)
	GetAddress(ctx context.Context, userID domain.UserId) (domain.Address, error)
}

type service struct {
	repo infrastructure.Repository
	enc  encryption.Service
}

func NewService(repo infrastructure.Repository, enc encryption.Service) Service {
	return &service{repo: repo, enc: enc}
}
