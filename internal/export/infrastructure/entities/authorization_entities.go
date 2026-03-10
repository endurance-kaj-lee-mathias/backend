package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type AuthorizationRuleExportEntity struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	ViewerID  uuid.UUID
	Resource  string
	Effect    string
	CreatedAt time.Time
}
