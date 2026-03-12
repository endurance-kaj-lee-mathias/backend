package models

import "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/application"

type VeteranMoodResponse struct {
	VeteranID string              `json:"veteranId"`
	FirstName string              `json:"firstName"`
	LastName  string              `json:"lastName"`
	Image     string              `json:"image"`
	Entries   []MoodEntryResponse `json:"entries"`
}

func ToVeteranMoodResponseList(summaries []application.VeteranMoodSummary) []VeteranMoodResponse {
	result := make([]VeteranMoodResponse, 0, len(summaries))

	for _, s := range summaries {
		result = append(result, VeteranMoodResponse{
			VeteranID: s.VeteranID.String(),
			FirstName: s.FirstName,
			LastName:  s.LastName,
			Image:     s.Image,
			Entries:   ToResponseList(s.Entries),
		})
	}

	return result
}
