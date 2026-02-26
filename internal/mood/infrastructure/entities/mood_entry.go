package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
)

type MoodEntryEntity struct {
	ID             uuid.UUID `db:"id"`
	UserID         uuid.UUID `db:"user_id"`
	Date           time.Time `db:"date"`
	MoodScore      int       `db:"mood_score"`
	EncryptedNotes []byte    `db:"encrypted_notes"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func ToEntity(entry domain.MoodEntry, enc encryption.Service, userKey []byte) (MoodEntryEntity, error) {
	var encryptedNotes []byte

	if entry.Notes != nil {
		encrypted, err := enc.Encrypt([]byte(*entry.Notes), userKey)
		if err != nil {
			return MoodEntryEntity{}, err
		}
		encryptedNotes = encrypted
	}

	return MoodEntryEntity{
		ID:             entry.ID.UUID,
		UserID:         entry.UserID.UUID,
		Date:           entry.Date,
		MoodScore:      entry.MoodScore,
		EncryptedNotes: encryptedNotes,
		CreatedAt:      entry.CreatedAt,
		UpdatedAt:      entry.UpdatedAt,
	}, nil
}
