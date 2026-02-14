package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

type Repository interface {
	Create(ctx context.Context, ent entities.UserEntity) error
	FindByID(ctx context.Context, id uuid.UUID) (entities.UserEntity, error)
	FindByEmail(ctx context.Context, email string) (entities.UserEntity, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}
