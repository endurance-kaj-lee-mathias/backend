package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type JournalReport struct {
	VeteranID    uuid.UUID
	UserProfile  *UserProfileSection
	StressScores *ScoresPage
	MoodEntries  *MoodPage
}

type UserProfileSection struct {
	FirstName    string
	LastName     string
	Username     string
	About        string
	Introduction string
	Image        string
	PhoneNumber  *string
	IsPrivate    bool
}

type ScoresPage struct {
	Items []StressScoreItem
	Total int
}

type MoodPage struct {
	Items []MoodEntryItem
	Total int
}

type StressScoreItem struct {
	ID           uuid.UUID
	Score        float64
	Category     string
	ModelVersion string
	ComputedAt   time.Time
}

type MoodEntryItem struct {
	ID        uuid.UUID
	Date      time.Time
	MoodScore int
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
