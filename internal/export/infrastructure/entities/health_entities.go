package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type StressSampleExportEntity struct {
	ID                      uuid.UUID
	TimestampUTC            time.Time
	WindowMinutes           int
	EncryptedMeanHR         []byte
	EncryptedRMSSDms        []byte
	EncryptedRestingHR      []byte
	EncryptedSteps          []byte
	EncryptedSleepDebtHours []byte
	CreatedAt               time.Time
}

type StressScoreExportEntity struct {
	ID           uuid.UUID
	Score        float64
	Category     string
	ModelVersion string
	ComputedAt   time.Time
}

type MoodEntryExportEntity struct {
	ID             uuid.UUID
	Date           time.Time
	MoodScore      int
	EncryptedNotes []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
