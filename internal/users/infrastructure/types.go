package infrastructure

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

type Repository interface {
	Save(ctx context.Context, e entities.UserEntity) error
	FindByID(ctx context.Context, id uuid.UUID) (entities.UserEntity, error)
	FindByEmail(ctx context.Context, email string) (entities.UserEntity, error)

	AddSupportMember(ctx context.Context, veteranID, supportID uuid.UUID) error
	ListSupportMembers(ctx context.Context, veteranID uuid.UUID) ([]entities.UserEntity, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}
