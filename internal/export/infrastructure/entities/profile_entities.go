package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type UserExportEntity struct {
	ID                    uuid.UUID
	EncryptedEmail        []byte
	EncryptedUsername     []byte
	EncryptedFirstName    []byte
	EncryptedLastName     []byte
	EncryptedPhoneNumber  []byte
	EncryptedRoles        []byte
	EncryptedAbout        []byte
	EncryptedIntroduction []byte
	Image                 string
	IsPrivate             bool
	EncryptedUserKey      []byte
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type AddressExportEntity struct {
	ID                  uuid.UUID
	EncryptedStreet     []byte
	EncryptedLocality   []byte
	EncryptedRegion     []byte
	EncryptedPostalCode []byte
	EncryptedCountry    []byte
	CreatedAt           time.Time
}
