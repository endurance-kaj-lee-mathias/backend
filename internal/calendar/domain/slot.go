package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Slot struct {
	ID         SlotId
	ProviderID uuid.UUID
	StartTime  time.Time
	EndTime    time.Time
	IsUrgent   bool
	IsBooked   bool
	SeriesID   *uuid.UUID
	Title      *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type SlotWithProvider struct {
	Slot
	ProviderUsername  string
	ProviderImage     string
	ProviderFirstName string
	ProviderLastName  string
}
