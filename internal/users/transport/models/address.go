package models

import (
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type AddressModel struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	Street     string    `json:"street"`
	Locality   string    `json:"locality"`
	Region     string    `json:"region"`
	PostalCode string    `json:"postalCode"`
	Country    string    `json:"country"`
}

func ToAddressModel(addr domain.Address) AddressModel {
	return AddressModel{
		ID:         addr.ID.UUID,
		UserID:     addr.UserID.UUID,
		Street:     addr.Street,
		Locality:   addr.Locality,
		Region:     addr.Region,
		PostalCode: addr.PostalCode,
		Country:    addr.Country,
	}
}

type UpsertAddressModel struct {
	Street     string `json:"street"`
	Locality   string `json:"locality"`
	Region     string `json:"region"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

func (m *UpsertAddressModel) Validate() error {
	if m.Street == "" {
		return InvalidStreet
	}
	if m.Locality == "" {
		return InvalidLocality
	}
	if m.PostalCode == "" {
		return InvalidPostalCode
	}
	if m.Country == "" {
		return InvalidCountry
	}
	return nil
}
