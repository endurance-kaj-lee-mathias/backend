package domain

import "time"

type MoodEntry struct {
	ID        MoodId
	UserID    UserId
	Date      time.Time
	MoodScore int
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMoodEntry(userID UserId, date time.Time, moodScore int, notes *string) (MoodEntry, error) {
	if moodScore < 0 || moodScore > 10 {
		return MoodEntry{}, InvalidScore
	}

	if notes != nil && len(*notes) > 500 {
		return MoodEntry{}, NotesTooLong
	}

	id, err := NewMoodId()
	if err != nil {
		return MoodEntry{}, err
	}

	now := time.Now().UTC()

	return MoodEntry{
		ID:        id,
		UserID:    userID,
		Date:      date,
		MoodScore: moodScore,
		Notes:     notes,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
