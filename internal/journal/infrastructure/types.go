package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/infrastructure/entities"
)

type Repository interface {
	GetUserProfile(ctx context.Context, userID uuid.UUID) (entities.UserProfileEntity, error)
	GetStressScoresPaginated(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entities.StressScoreRow, int, error)
	GetMoodEntriesPaginated(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entities.MoodEntryRow, int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
