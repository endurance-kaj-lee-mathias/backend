package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

type JournalResponse struct {
	VeteranID    string                                           `json:"veteranId"`
	UserProfile  UserProfileResponse                              `json:"userProfile"`
	StressScores *response.PaginatedResponse[StressScoreResponse] `json:"stressScores,omitempty"`
	MoodEntries  *response.PaginatedResponse[MoodEntryResponse]   `json:"moodEntries,omitempty"`
}

type UserProfileResponse struct {
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Username     string  `json:"username"`
	About        string  `json:"about"`
	Introduction string  `json:"introduction"`
	Image        string  `json:"image"`
	PhoneNumber  *string `json:"phoneNumber,omitempty"`
	IsPrivate    bool    `json:"isPrivate"`
}

type StressScoreResponse struct {
	ID           uuid.UUID `json:"id"`
	Score        float64   `json:"score"`
	Category     string    `json:"category"`
	ModelVersion string    `json:"modelVersion"`
	ComputedAt   time.Time `json:"computedAt"`
}

type MoodEntryResponse struct {
	ID        uuid.UUID `json:"id"`
	Date      string    `json:"date"`
	MoodScore int       `json:"moodScore"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToJournalResponse(report domain.JournalReport, limit, offset int) (JournalResponse, error) {
	if report.UserProfile == nil {
		return JournalResponse{}, domain.MissingUserProfile
	}

	p := report.UserProfile
	jr := JournalResponse{
		VeteranID: report.VeteranID.String(),
		UserProfile: UserProfileResponse{
			FirstName:    p.FirstName,
			LastName:     p.LastName,
			Username:     p.Username,
			About:        p.About,
			Introduction: p.Introduction,
			Image:        p.Image,
			PhoneNumber:  p.PhoneNumber,
			IsPrivate:    p.IsPrivate,
		},
	}

	if report.StressScores != nil {
		items := make([]StressScoreResponse, 0, len(report.StressScores.Items))

		for _, item := range report.StressScores.Items {
			items = append(items, StressScoreResponse{
				ID:           item.ID,
				Score:        item.Score,
				Category:     item.Category,
				ModelVersion: item.ModelVersion,
				ComputedAt:   item.ComputedAt,
			})
		}

		p := response.NewPaginated(items, report.StressScores.Total, limit, offset)
		jr.StressScores = &p
	}

	if report.MoodEntries != nil {
		items := make([]MoodEntryResponse, 0, len(report.MoodEntries.Items))

		for _, item := range report.MoodEntries.Items {
			items = append(items, MoodEntryResponse{
				ID:        item.ID,
				Date:      item.Date.Format("2006-01-02"),
				MoodScore: item.MoodScore,
				Notes:     item.Notes,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			})
		}

		p := response.NewPaginated(items, report.MoodEntries.Total, limit, offset)
		jr.MoodEntries = &p
	}

	return jr, nil
}
