package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

type Repository interface {
	Create(ctx context.Context, ent entities.UserEntity) error
	FindByID(ctx context.Context, id uuid.UUID) (entities.UserEntity, error)
	FindByEmail(ctx context.Context, email string) (entities.UserEntity, error)
	FindByUsername(ctx context.Context, username string) (entities.UserEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePhoneNumber(ctx context.Context, id uuid.UUID, phoneNumber *string) error
	UpdateIntroduction(ctx context.Context, id uuid.UUID, encrypted []byte) error
	UpdateAbout(ctx context.Context, id uuid.UUID, encrypted []byte) error
	UpdateImage(ctx context.Context, id uuid.UUID, image string) error
	InsertAddress(ctx context.Context, ent entities.AddressEntity) error
	FindAddressByUserID(ctx context.Context, userID uuid.UUID) (entities.AddressEntity, error)
	GetEncryptedUserKey(ctx context.Context, userID uuid.UUID) ([]byte, error)
}

type repository struct {
	db  *sql.DB
	enc encryption.Service
}

func NewRepository(db *sql.DB, enc encryption.Service) Repository {
	return &repository{db: db, enc: enc}
}
