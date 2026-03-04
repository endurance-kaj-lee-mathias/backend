package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type StressScore struct {
	ID           ScoreId
	UserID       uuid.UUID
	Score        float64
	Category     string
	ModelVersion string
	ComputedAt   time.Time
}

func NewStressScore(userID uuid.UUID, score float64, category string, modelVersion string) (StressScore, error) {
	id, err := NewScoreId()
	if err != nil {
		return StressScore{}, err
	}

	return StressScore{
		ID:           id,
		UserID:       userID,
		Score:        score,
		Category:     category,
		ModelVersion: modelVersion,
		ComputedAt:   time.Now().UTC(),
	}, nil
}
