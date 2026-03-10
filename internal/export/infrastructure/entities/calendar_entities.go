package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type AppointmentExportEntity struct {
	ID                        uuid.UUID
	SlotID                    uuid.UUID
	VeteranID                 uuid.UUID
	ProviderID                uuid.UUID
	Status                    string
	StartTime                 time.Time
	EndTime                   time.Time
	IsUrgent                  bool
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	ProviderEncryptedUsername []byte
	ProviderEncryptedUserKey  []byte
}

type SlotExportEntity struct {
	ID        uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	IsUrgent  bool
	IsBooked  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
