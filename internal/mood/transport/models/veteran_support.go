package models

import (
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
)

type VeteranSupportResponse struct {
	ID            string     `json:"id"`
	FirstName     string     `json:"firstName"`
	LastName      string     `json:"lastName"`
	Image         string     `json:"image"`
	LastUpdatedAt *time.Time `json:"lastUpdatedAt,omitempty"`
	LatestScore   *int       `json:"latestScore,omitempty"`
}

func ToVeteranSupportResponseList(summaries []domain.VeteranMoodSummary) []VeteranSupportResponse {
	result := make([]VeteranSupportResponse, 0, len(summaries))

	for _, s := range summaries {
		result = append(result, VeteranSupportResponse{
			ID:            s.VeteranID.String(),
			FirstName:     s.FirstName,
			LastName:      s.LastName,
			Image:         s.Image,
			LastUpdatedAt: s.LastUpdatedAt,
			LatestScore:   s.LatestScore,
		})
	}

	return result
}
