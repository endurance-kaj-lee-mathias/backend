package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
)

type StressScoreResponse struct {
	ID           uuid.UUID `json:"id"`
	Score        float64   `json:"score"`
	Category     string    `json:"category"`
	ModelVersion string    `json:"modelVersion"`
	ComputedAt   time.Time `json:"computedAt"`
}

func ToStressScoreResponse(score domain.StressScore) StressScoreResponse {
	return StressScoreResponse{
		ID:           score.ID.UUID,
		Score:        score.Score,
		Category:     score.Category,
		ModelVersion: score.ModelVersion,
		ComputedAt:   score.ComputedAt,
	}
}

func ToStressScoreResponseList(scores []domain.StressScore) []StressScoreResponse {
	result := make([]StressScoreResponse, 0, len(scores))

	for _, score := range scores {
		result = append(result, ToStressScoreResponse(score))
	}

	return result
}
