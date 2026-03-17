package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type CalendarEventEntity struct {
	ID                         uuid.UUID
	ProviderID                 uuid.UUID
	VeteranID                  uuid.UUID
	StartTime                  time.Time
	EndTime                    time.Time
	UpdatedAt                  time.Time
	ProviderEncryptedFirstName []byte
	ProviderEncryptedLastName  []byte
	ProviderEncryptedUserKey   []byte
	VeteranEncryptedFirstName  []byte
	VeteranEncryptedLastName   []byte
	VeteranEncryptedUserKey    []byte
}
