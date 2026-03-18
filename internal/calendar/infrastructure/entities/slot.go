package entities

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

type SlotEntity struct {
	ID               uuid.UUID      `db:"id"`
	ProviderID       uuid.UUID      `db:"provider_id"`
	StartTime        time.Time      `db:"start_time"`
	EndTime          time.Time      `db:"end_time"`
	IsUrgent         bool           `db:"is_urgent"`
	IsBooked         bool           `db:"is_booked"`
	SeriesID         *uuid.UUID     `db:"series_id"`
	AppointmentTitle sql.NullString `db:"title"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}

func SlotToEntity(slot domain.Slot) SlotEntity {
	var title sql.NullString
	if slot.Title != nil {
		title = sql.NullString{String: *slot.Title, Valid: true}
	}

	return SlotEntity{
		ID:               slot.ID.UUID,
		ProviderID:       slot.ProviderID,
		StartTime:        slot.StartTime,
		EndTime:          slot.EndTime,
		IsUrgent:         slot.IsUrgent,
		IsBooked:         slot.IsBooked,
		SeriesID:         slot.SeriesID,
		AppointmentTitle: title,
		CreatedAt:        slot.CreatedAt,
		UpdatedAt:        slot.UpdatedAt,
	}
}

func SlotFromEntity(ent SlotEntity) domain.Slot {
	var title *string
	if ent.AppointmentTitle.Valid {
		title = &ent.AppointmentTitle.String
	}

	return domain.Slot{
		ID:         domain.SlotId{UUID: ent.ID},
		ProviderID: ent.ProviderID,
		StartTime:  ent.StartTime,
		EndTime:    ent.EndTime,
		IsUrgent:   ent.IsUrgent,
		IsBooked:   ent.IsBooked,
		SeriesID:   ent.SeriesID,
		Title:      title,
		CreatedAt:  ent.CreatedAt,
		UpdatedAt:  ent.UpdatedAt,
	}
}

func SlotsFromEntities(ents []SlotEntity) []domain.Slot {
	slots := make([]domain.Slot, 0, len(ents))
	for _, ent := range ents {
		slots = append(slots, SlotFromEntity(ent))
	}
	return slots
}
