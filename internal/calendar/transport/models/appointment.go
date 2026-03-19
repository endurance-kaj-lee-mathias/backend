package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

type AppointmentResponse struct {
	ID                uuid.UUID `json:"id"`
	SlotID            uuid.UUID `json:"slotId"`
	VeteranID         uuid.UUID `json:"veteranId"`
	Title             *string   `json:"title,omitempty"`
	Status            string    `json:"status"`
	StartTime         time.Time `json:"startTime"`
	EndTime           time.Time `json:"endTime"`
	ProviderUsername  string    `json:"providerUsername"`
	ProviderImage     string    `json:"providerImage"`
	ProviderFirstName string    `json:"providerFirstName"`
	ProviderLastName  string    `json:"providerLastName"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func ToAppointmentModel(a domain.Appointment) AppointmentResponse {
	return AppointmentResponse{
		ID:                a.ID.UUID,
		SlotID:            a.SlotID,
		VeteranID:         a.VeteranID,
		Title:             a.Title,
		Status:            string(a.Status),
		StartTime:         a.StartTime,
		EndTime:           a.EndTime,
		ProviderUsername:  a.ProviderUsername,
		ProviderImage:     a.ProviderImage,
		ProviderFirstName: a.ProviderFirstName,
		ProviderLastName:  a.ProviderLastName,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

func ToAppointmentModels(appointments []domain.Appointment) []AppointmentResponse {
	out := make([]AppointmentResponse, 0, len(appointments))
	for _, a := range appointments {
		out = append(out, ToAppointmentModel(a))
	}
	return out
}

type EventResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToEventModels(events []domain.Event) []EventResponse {
	out := make([]EventResponse, 0, len(events))
	for _, e := range events {
		out = append(out, EventResponse{
			ID:        e.ID,
			Title:     e.Title,
			StartTime: e.StartTime,
			EndTime:   e.EndTime,
			UpdatedAt: e.UpdatedAt,
		})
	}
	return out
}
