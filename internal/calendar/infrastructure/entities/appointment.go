package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

const StatusCancelled = "CANCELLED"

type AppointmentEntity struct {
	ID        uuid.UUID `db:"id"`
	SlotID    uuid.UUID `db:"slot_id"`
	VeteranID uuid.UUID `db:"veteran_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AppointmentWithSlotEntity struct {
	AppointmentEntity
	SlotProviderID uuid.UUID `db:"provider_id"`
}

func AppointmentToEntity(a domain.Appointment) AppointmentEntity {
	return AppointmentEntity{
		ID:        a.ID.UUID,
		SlotID:    a.SlotID,
		VeteranID: a.VeteranID,
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func AppointmentFromEntity(ent AppointmentEntity) domain.Appointment {
	return domain.Appointment{
		ID:        domain.AppointmentId{UUID: ent.ID},
		SlotID:    ent.SlotID,
		VeteranID: ent.VeteranID,
		Status:    domain.AppointmentStatus(ent.Status),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
