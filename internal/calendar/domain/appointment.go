package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type AppointmentStatus string

const (
	StatusBooked    AppointmentStatus = "BOOKED"
	StatusCancelled AppointmentStatus = "CANCELLED"
	StatusCompleted AppointmentStatus = "COMPLETED"
)

type Appointment struct {
	ID        AppointmentId
	SlotID    uuid.UUID
	VeteranID uuid.UUID
	Title     *string
	Status    AppointmentStatus
	StartTime time.Time
	EndTime   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
