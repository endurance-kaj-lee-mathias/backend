package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
)

type Repository interface {
	Create(ctx context.Context, veteranID, memberId uuid.UUID) (entities.MemberEntity, error)
	ReadAll(ctx context.Context, id uuid.UUID) ([]entities.MemberEntity, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}
