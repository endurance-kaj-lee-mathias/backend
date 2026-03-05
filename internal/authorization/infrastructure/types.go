package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
)

type Repository interface {
	Create(ctx context.Context, rule domain.Rule) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Rule, error)
	FindByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.Rule, error)
	FindRule(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID, resource string) (*domain.Rule, error)
	DeleteByOwnerAndViewer(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID) error
	GetPrivacy(ctx context.Context, userID uuid.UUID) (bool, error)
	HasSupportRelationship(ctx context.Context, userA uuid.UUID, userB uuid.UUID) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
