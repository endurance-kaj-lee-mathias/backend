package entities

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

const StatusCancelled = "CANCELLED"

type AppointmentEntity struct {
	ID        uuid.UUID      `db:"id"`
	SlotID    uuid.UUID      `db:"slot_id"`
	VeteranID uuid.UUID      `db:"veteran_id"`
	Title     sql.NullString `db:"title"`
	Status    string         `db:"status"`
	StartTime time.Time      `db:"start_time"`
	EndTime   time.Time      `db:"end_time"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type AppointmentWithSlotEntity struct {
	AppointmentEntity
	SlotProviderID uuid.UUID `db:"provider_id"`
}

func AppointmentToEntity(a domain.Appointment) AppointmentEntity {
	var title sql.NullString
	if a.Title != nil {
		title = sql.NullString{String: *a.Title, Valid: true}
	}

	return AppointmentEntity{
		ID:        a.ID.UUID,
		SlotID:    a.SlotID,
		VeteranID: a.VeteranID,
		Title:     title,
		Status:    string(a.Status),
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func AppointmentFromEntity(ent AppointmentEntity) domain.Appointment {
	var title *string
	if ent.Title.Valid {
		title = &ent.Title.String
	}

	return domain.Appointment{
		ID:        domain.AppointmentId{UUID: ent.ID},
		SlotID:    ent.SlotID,
		VeteranID: ent.VeteranID,
		StartTime: ent.StartTime,
		EndTime:   ent.EndTime,
		Title:     title,
		Status:    domain.AppointmentStatus(ent.Status),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
