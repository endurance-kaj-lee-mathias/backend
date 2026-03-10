package domain

import (
	"time"
)

type CalendarData struct {
	Appointments []AppointmentData `json:"appointments"`
	Slots        []SlotData        `json:"slots"`
}

type AppointmentData struct {
	ID               string    `json:"id"`
	SlotID           string    `json:"slotId"`
	VeteranID        string    `json:"veteranId"`
	ProviderID       string    `json:"providerId"`
	ProviderUsername string    `json:"providerUsername"`
	Status           string    `json:"status"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	IsUrgent         bool      `json:"isUrgent"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type SlotData struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	IsUrgent  bool      `json:"isUrgent"`
	IsBooked  bool      `json:"isBooked"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
