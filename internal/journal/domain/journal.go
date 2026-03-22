package domain

import "github.com/gofrs/uuid"

type JournalReport struct {
	VeteranID   uuid.UUID
	UserProfile *UserProfileSection
	Weekly      *WeeklyPage
}

type UserProfileSection struct {
	FirstName    string
	LastName     string
	Username     string
	About        string
	Introduction string
	Image        string
	PhoneNumber  *string
	IsPrivate    bool
	RiskLevel    string
}

type WeeklyPage struct {
	Days  []DailyAverage
	Total int
}

type DailyAverage struct {
	Date      string
	AvgMood   float64
	AvgStress *float64
	Notes     []string
}
