package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type AddressEntity struct {
	ID                  uuid.UUID `db:"id"`
	UserID              uuid.UUID `db:"user_id"`
	EncryptedStreet     []byte    `db:"encrypted_street"`
	EncryptedLocality   []byte    `db:"encrypted_locality"`
	EncryptedRegion     []byte    `db:"encrypted_region"`
	EncryptedPostalCode []byte    `db:"encrypted_postal_code"`
	EncryptedCountry    []byte    `db:"encrypted_country"`
	CreatedAt           time.Time `db:"created_at"`
}

func AddressFromEntity(ent AddressEntity, userKey []byte, enc encryption.Service) (domain.Address, error) {
	streetBytes, err := enc.Decrypt(ent.EncryptedStreet, userKey)
	if err != nil {
		return domain.Address{}, err
	}

	localityBytes, err := enc.Decrypt(ent.EncryptedLocality, userKey)
	if err != nil {
		return domain.Address{}, err
	}

	regionBytes, err := enc.Decrypt(ent.EncryptedRegion, userKey)
	if err != nil {
		return domain.Address{}, err
	}

	postalCodeBytes, err := enc.Decrypt(ent.EncryptedPostalCode, userKey)
	if err != nil {
		return domain.Address{}, err
	}

	countryBytes, err := enc.Decrypt(ent.EncryptedCountry, userKey)
	if err != nil {
		return domain.Address{}, err
	}

	id, err := domain.ParseAddressId(ent.ID.String())
	if err != nil {
		return domain.Address{}, err
	}

	userID, err := domain.ParseId(ent.UserID.String())
	if err != nil {
		return domain.Address{}, err
	}

	return domain.Address{
		ID:         id,
		UserID:     userID,
		Street:     string(streetBytes),
		Locality:   string(localityBytes),
		Region:     string(regionBytes),
		PostalCode: string(postalCodeBytes),
		Country:    string(countryBytes),
		CreatedAt:  ent.CreatedAt,
	}, nil
}

func AddressToEntity(a domain.Address, userKey []byte, enc encryption.Service) (AddressEntity, error) {
	encStreet, err := enc.Encrypt([]byte(a.Street), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	encLocality, err := enc.Encrypt([]byte(a.Locality), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	encRegion, err := enc.Encrypt([]byte(a.Region), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	encPostalCode, err := enc.Encrypt([]byte(a.PostalCode), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	encCountry, err := enc.Encrypt([]byte(a.Country), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	return AddressEntity{
		ID:                  a.ID.UUID,
		UserID:              a.UserID.UUID,
		EncryptedStreet:     encStreet,
		EncryptedLocality:   encLocality,
		EncryptedRegion:     encRegion,
		EncryptedPostalCode: encPostalCode,
		EncryptedCountry:    encCountry,
		CreatedAt:           a.CreatedAt,
	}, nil
}
