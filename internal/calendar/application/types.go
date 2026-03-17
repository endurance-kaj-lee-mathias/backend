package application

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
)

type Service interface {
	CreateSlot(ctx context.Context, providerID uuid.UUID, roles []string, startTime, endTime time.Time, isUrgent bool, isRecurring bool) (domain.Slot, error)
	GetSlots(ctx context.Context, from, to time.Time, providerID *uuid.UUID) ([]domain.Slot, error)
	DeleteSlot(ctx context.Context, userID uuid.UUID, roles []string, slotID domain.SlotId) error
	DeleteSlotsBySeries(ctx context.Context, providerID uuid.UUID, seriesID uuid.UUID) error
	BookSlot(ctx context.Context, veteranID uuid.UUID, roles []string, slotID domain.SlotId, urgent bool) (domain.Appointment, error)
	CancelAppointment(ctx context.Context, userID uuid.UUID, appointmentID domain.AppointmentId) error
	DeleteMySlots(ctx context.Context, providerID uuid.UUID) error
	GetCalendarEvents(ctx context.Context, userID uuid.UUID) ([]domain.Event, error)
}

type service struct {
	repo             infrastructure.Repository
	enc              encryption.Service
	minUrgentMinutes int
}

func NewService(repo infrastructure.Repository, enc encryption.Service, minUrgentMinutes int) Service {
	return &service{repo: repo, enc: enc, minUrgentMinutes: minUrgentMinutes}
}
