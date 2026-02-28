package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure/entities"
)

type Repository interface {
	Upsert(ctx context.Context, ent entities.MoodEntryEntity) error
	FindVeteransWithoutEntryInLast24Hours(ctx context.Context, veteranRoleHash string) ([]uuid.UUID, error)
	FindDeviceTokensByUserID(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type UserKeyReader interface {
	GetEncryptedUserKey(ctx context.Context, userID domain.UserId) ([]byte, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

type userKeyReader struct {
	db  *sql.DB
	enc encryption.Service
}

func NewUserKeyReader(db *sql.DB, enc encryption.Service) UserKeyReader {
	return &userKeyReader{db: db, enc: enc}
}
