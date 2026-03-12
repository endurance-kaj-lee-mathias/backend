package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type UserProfileEntity struct {
	ID                    uuid.UUID
	EncryptedUserKey      []byte
	EncryptedFirstName    []byte
	EncryptedLastName     []byte
	EncryptedUsername     []byte
	EncryptedAbout        []byte
	EncryptedIntroduction []byte
	EncryptedPhoneNumber  []byte
	Image                 string
	IsPrivate             bool
}

type StressScoreRow struct {
	ID           uuid.UUID
	Score        float64
	Category     string
	ModelVersion string
	ComputedAt   time.Time
}

type MoodEntryRow struct {
	ID             uuid.UUID
	Date           time.Time
	MoodScore      int
	EncryptedNotes []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

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
