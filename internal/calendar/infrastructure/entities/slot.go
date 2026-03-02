package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

type SlotEntity struct {
	ID         uuid.UUID `db:"id"`
	ProviderID uuid.UUID `db:"provider_id"`
	StartTime  time.Time `db:"start_time"`
	EndTime    time.Time `db:"end_time"`
	IsUrgent   bool      `db:"is_urgent"`
	IsBooked   bool      `db:"is_booked"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func SlotToEntity(slot domain.Slot) SlotEntity {
	return SlotEntity{
		ID:         slot.ID.UUID,
		ProviderID: slot.ProviderID,
		StartTime:  slot.StartTime,
		EndTime:    slot.EndTime,
		IsUrgent:   slot.IsUrgent,
		IsBooked:   slot.IsBooked,
		CreatedAt:  slot.CreatedAt,
		UpdatedAt:  slot.UpdatedAt,
	}
}

func SlotFromEntity(ent SlotEntity) domain.Slot {
	return domain.Slot{
		ID:         domain.SlotId{UUID: ent.ID},
		ProviderID: ent.ProviderID,
		StartTime:  ent.StartTime,
		EndTime:    ent.EndTime,
		IsUrgent:   ent.IsUrgent,
		IsBooked:   ent.IsBooked,
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
