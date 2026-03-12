package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Repository interface {
	GetUserProfile(ctx context.Context, userID uuid.UUID) (UserProfileEntity, error)
	GetStressScoresPaginated(ctx context.Context, userID uuid.UUID, limit, offset int) ([]StressScoreRow, int, error)
	GetMoodEntriesPaginated(ctx context.Context, userID uuid.UUID, limit, offset int) ([]MoodEntryRow, int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
