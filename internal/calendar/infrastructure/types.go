package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure/entities"
)

type Repository interface {
	CreateSlot(ctx context.Context, ent entities.SlotEntity) error
	GetSlotsByRange(ctx context.Context, from, to time.Time, providerID *uuid.UUID) ([]entities.SlotEntity, error)
	GetSlotByID(ctx context.Context, id uuid.UUID) (entities.SlotEntity, error)
	DeleteSlot(ctx context.Context, id uuid.UUID) error
	AtomicBookSlot(ctx context.Context, id uuid.UUID, now time.Time) (int64, error)
	CreateAppointment(ctx context.Context, ent entities.AppointmentEntity) error
	GetAppointmentWithSlot(ctx context.Context, id uuid.UUID) (entities.AppointmentWithSlotEntity, error)
	CancelAppointment(ctx context.Context, appointmentID uuid.UUID, now time.Time) error
	CheckSlotOverlap(ctx context.Context, providerID uuid.UUID, start, end time.Time) (bool, error)
	GetUrgentSlotMinutesForDate(ctx context.Context, providerID uuid.UUID, date time.Time) (int, error)
	DeleteSlotsByProviderID(ctx context.Context, providerID uuid.UUID) error
	DeleteFutureSlotsBySeries(ctx context.Context, seriesID uuid.UUID, providerID uuid.UUID) error
	GetEventsByUserID(ctx context.Context, userID uuid.UUID) ([]entities.CalendarEventEntity, error)
	GetAppointmentsByDay(ctx context.Context, veteranID uuid.UUID, dayStart, dayEnd time.Time) ([]entities.AppointmentWithSlotEntity, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
