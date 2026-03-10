package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type InviteExportEntity struct {
	ID                      uuid.UUID
	OtherUserID             uuid.UUID
	OtherEncryptedUsername  []byte
	OtherEncryptedFirstName []byte
	OtherEncryptedLastName  []byte
	OtherEncryptedUserKey   []byte
	OtherImage              *string
	Status                  string
	Note                    *string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
