package models

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/domain"
)

type JournalResponse struct {
	VeteranID   string               `json:"veteranId"`
	UserProfile *UserProfileResponse `json:"profile,omitempty"`
	Weekly      *WeeklyResponse      `json:"weekly,omitempty"`
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

type WeeklyResponse struct {
	Days       []DailyAverageResponse `json:"days"`
	TotalWeeks int                    `json:"totalWeeks"`
	Week       int                    `json:"week"`
}

type DailyAverageResponse struct {
	Date      string   `json:"date"`
	AvgMood   float64  `json:"avgMood"`
	AvgStress *float64 `json:"avgStress"`
}

func ToJournalResponse(report domain.JournalReport, weekOffset int) (JournalResponse, error) {
	if report.UserProfile == nil && report.Weekly == nil {
		return JournalResponse{VeteranID: report.VeteranID.String()}, nil
	}

	jr := JournalResponse{
		VeteranID: report.VeteranID.String(),
	}

	if report.UserProfile != nil {
		p := report.UserProfile
		prof := UserProfileResponse{
			FirstName:    p.FirstName,
			LastName:     p.LastName,
			Username:     p.Username,
			About:        p.About,
			Introduction: p.Introduction,
			Image:        p.Image,
			PhoneNumber:  p.PhoneNumber,
			IsPrivate:    p.IsPrivate,
		}
		jr.UserProfile = &prof
	}

	if report.Weekly != nil {
		days := make([]DailyAverageResponse, 0, len(report.Weekly.Days))

		for _, d := range report.Weekly.Days {
			days = append(days, DailyAverageResponse{
				Date:      d.Date,
				AvgMood:   d.AvgMood,
				AvgStress: d.AvgStress,
			})
		}

		jr.Weekly = &WeeklyResponse{
			Days:       days,
			TotalWeeks: report.Weekly.Total,
			Week:       weekOffset,
		}
	}

	return jr, nil
}
