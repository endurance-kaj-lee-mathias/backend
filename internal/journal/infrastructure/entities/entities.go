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
	EncryptedRoles        []byte
	Image                 string
	IsPrivate             bool
	RiskLevel             string
}

type DailyAverageRow struct {
	Date      time.Time
	AvgMood   float64
	AvgStress *float64
	Total     int
}

type MoodEntryNoteRow struct {
	Date           time.Time
	EncryptedNotes []byte
}
