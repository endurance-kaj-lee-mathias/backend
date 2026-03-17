package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

type SlotResponse struct {
	ID         uuid.UUID  `json:"id"`
	ProviderID uuid.UUID  `json:"providerId"`
	StartTime  time.Time  `json:"startTime"`
	EndTime    time.Time  `json:"endTime"`
	IsUrgent   bool       `json:"isUrgent"`
	IsBooked   bool       `json:"isBooked"`
	SeriesID   *uuid.UUID `json:"seriesId,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func ToSlotModel(s domain.Slot) SlotResponse {
	return SlotResponse{
		ID:         s.ID.UUID,
		ProviderID: s.ProviderID,
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
		IsUrgent:   s.IsUrgent,
		IsBooked:   s.IsBooked,
		SeriesID:   s.SeriesID,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

func ToSlotModels(slots []domain.Slot) []SlotResponse {
	out := make([]SlotResponse, 0, len(slots))
	for _, s := range slots {
		out = append(out, ToSlotModel(s))
	}
	return out
}
