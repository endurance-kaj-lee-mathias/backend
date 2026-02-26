package models

type MoodEntryRequest struct {
	Date      string  `json:"date"`
	MoodScore int     `json:"moodScore"`
	Notes     *string `json:"notes,omitempty"`
}

func (m *MoodEntryRequest) Validate() error {
	if m.Date == "" {
		return InvalidDate
	}
	if m.MoodScore < 1 || m.MoodScore > 10 {
		return InvalidMoodScore
	}
	if m.Notes != nil && len(*m.Notes) > 500 {
		return NotesTooLong
	}
	return nil
}
