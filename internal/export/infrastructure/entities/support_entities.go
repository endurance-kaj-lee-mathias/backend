package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type SupportMemberExportEntity struct {
	ID                uuid.UUID
	EncryptedEmail    []byte
	EncryptedUsername []byte
	EncryptedFirst    []byte
	EncryptedLast     []byte
	Image             string
	EncryptedUserKey  []byte
	CreatedAt         time.Time
}
