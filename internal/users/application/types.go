package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/keycloak"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"
)

type Service interface {
	GetOrCreate(ctx context.Context, id domain.UserId, email string, username string, firstName string, lastName string, phoneNumber string, street string, locality string, region string, postalCode string, country string, roles []domain.Role) (domain.User, error)
	GetByID(ctx context.Context, id domain.UserId) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	DeleteUser(ctx context.Context, id domain.UserId) error
	UpdatePhoneNumber(ctx context.Context, id domain.UserId, phoneNumber *string) error
	UpdateFirstName(ctx context.Context, id domain.UserId, firstName string) error
	UpdateLastName(ctx context.Context, id domain.UserId, lastName string) error
	UpdateIntroduction(ctx context.Context, id domain.UserId, introduction string) error
	UpdateAbout(ctx context.Context, id domain.UserId, about string) error
	UpdateImage(ctx context.Context, id domain.UserId, image string) error
	UpdatePrivacy(ctx context.Context, id domain.UserId, isPrivate bool) error
	UpdateRiskLevel(ctx context.Context, id domain.UserId, riskLevel domain.RiskLevel) error
	UpsertAddress(ctx context.Context, userID domain.UserId, street string, locality string, region string, postalCode string, country string) (domain.Address, error)
	GetAddress(ctx context.Context, userID domain.UserId) (domain.Address, error)
	UpsertDevice(ctx context.Context, userID domain.UserId, deviceToken string, platform string) error
	DeleteDevice(ctx context.Context, deviceToken string) error
	FindDeviceTokensByUserID(ctx context.Context, userID domain.UserId) ([]string, error)
	AssignRole(ctx context.Context, userID domain.UserId, roleName string) error
}

type service struct {
	repo infrastructure.Repository
	enc  encryption.Service
	kc   keycloak.Client
}

func NewService(repo infrastructure.Repository, enc encryption.Service, kc keycloak.Client) Service {
	return &service{repo: repo, enc: enc, kc: kc}
}
