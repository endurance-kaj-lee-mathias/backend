package domain

import (
	"time"
)

type HealthData struct {
	StressSamples []StressSampleData `json:"stressSamples"`
	StressScores  []StressScoreData  `json:"stressScores"`
	MoodEntries   []MoodEntryData    `json:"moodEntries"`
}

type StressSampleData struct {
	ID             string    `json:"id"`
	TimestampUTC   time.Time `json:"timestampUtc"`
	WindowMinutes  int       `json:"windowMinutes"`
	MeanHR         float64   `json:"meanHr"`
	RMSSDms        float64   `json:"rmssdMs"`
	RestingHR      *float64  `json:"restingHr,omitempty"`
	Steps          *int      `json:"steps,omitempty"`
	SleepDebtHours *float64  `json:"sleepDebtHours,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type StressScoreData struct {
	ID           string    `json:"id"`
	Score        float64   `json:"score"`
	Category     string    `json:"category"`
	ModelVersion string    `json:"modelVersion"`
	ComputedAt   time.Time `json:"computedAt"`
}

type MoodEntryData struct {
	ID        string    `json:"id"`
	Date      time.Time `json:"date"`
	MoodScore int       `json:"moodScore"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
