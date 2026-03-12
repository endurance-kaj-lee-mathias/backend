package entities

import (
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
