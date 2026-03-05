package models

import (
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
)

type MoodEntryRequest struct {
	Date      string  `json:"date"`
	MoodScore int     `json:"moodScore"`
	Notes     *string `json:"notes,omitempty"`
}

func (m *MoodEntryRequest) Validate() error {
	if m.Date == "" {
		return InvalidDate
	}
	if m.MoodScore < 0 || m.MoodScore > 10 {
		return InvalidMoodScore
	}
	if m.Notes != nil && len(*m.Notes) > 500 {
		return NotesTooLong
	}
	return nil
}

type MoodEntryResponse struct {
	ID        string    `json:"id"`
	Date      string    `json:"date"`
	MoodScore int       `json:"moodScore"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToResponse(e domain.MoodEntry) MoodEntryResponse {
	return MoodEntryResponse{
		ID:        e.ID.String(),
		Date:      e.Date.Format("2006-01-02"),
		MoodScore: e.MoodScore,
		Notes:     e.Notes,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func ToResponseList(entries []domain.MoodEntry) []MoodEntryResponse {
	result := make([]MoodEntryResponse, 0, len(entries))

	for _, e := range entries {
		result = append(result, MoodEntryResponse{
			ID:        e.ID.String(),
			Date:      e.Date.Format("2006-01-02"),
			MoodScore: e.MoodScore,
			Notes:     e.Notes,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	return result
}
