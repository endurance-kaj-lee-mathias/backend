package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

type AppointmentResponse struct {
	ID        uuid.UUID `json:"id"`
	SlotID    uuid.UUID `json:"slotId"`
	VeteranID uuid.UUID `json:"veteranId"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToAppointmentModel(a domain.Appointment) AppointmentResponse {
	return AppointmentResponse{
		ID:        a.ID.UUID,
		SlotID:    a.SlotID,
		VeteranID: a.VeteranID,
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
