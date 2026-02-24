package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type AddressEntity struct {
	ID                   uuid.UUID `db:"id"`
	UserID               uuid.UUID `db:"user_id"`
	EncryptedStreet      []byte    `db:"encrypted_street"`
	EncryptedHouseNumber []byte    `db:"encrypted_house_number"`
	EncryptedPostalCode  []byte    `db:"encrypted_postal_code"`
	EncryptedCity        []byte    `db:"encrypted_city"`
	EncryptedCountry     []byte    `db:"encrypted_country"`
	CreatedAt            time.Time `db:"created_at"`
}

func AddressFromEntity(ent AddressEntity, userKey []byte, enc encryption.Service) (domain.Address, error) {
	streetBytes, err := enc.Decrypt(ent.EncryptedStreet, userKey)
	if err != nil {
		return domain.Address{}, err
	}

	houseNumberBytes, err := enc.Decrypt(ent.EncryptedHouseNumber, userKey)
	if err != nil {
		return domain.Address{}, err
	}

	postalCodeBytes, err := enc.Decrypt(ent.EncryptedPostalCode, userKey)
	if err != nil {
		return domain.Address{}, err
	}

	cityBytes, err := enc.Decrypt(ent.EncryptedCity, userKey)
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
		ID:          id,
		UserID:      userID,
		Street:      string(streetBytes),
		HouseNumber: string(houseNumberBytes),
		PostalCode:  string(postalCodeBytes),
		City:        string(cityBytes),
		Country:     string(countryBytes),
		CreatedAt:   ent.CreatedAt,
	}, nil
}

func AddressToEntity(a domain.Address, userKey []byte, enc encryption.Service) (AddressEntity, error) {
	encStreet, err := enc.Encrypt([]byte(a.Street), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	encHouseNumber, err := enc.Encrypt([]byte(a.HouseNumber), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	encPostalCode, err := enc.Encrypt([]byte(a.PostalCode), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	encCity, err := enc.Encrypt([]byte(a.City), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	encCountry, err := enc.Encrypt([]byte(a.Country), userKey)
	if err != nil {
		return AddressEntity{}, err
	}

	return AddressEntity{
		ID:                   a.ID.UUID,
		UserID:               a.UserID.UUID,
		EncryptedStreet:      encStreet,
		EncryptedHouseNumber: encHouseNumber,
		EncryptedPostalCode:  encPostalCode,
		EncryptedCity:        encCity,
		EncryptedCountry:     encCountry,
		CreatedAt:            a.CreatedAt,
	}, nil
}
