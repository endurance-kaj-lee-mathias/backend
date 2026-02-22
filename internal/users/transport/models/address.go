package models

import "github.com/gofrs/uuid"

type UpsertAddressModel struct {
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	PostalCode  string `json:"postalCode"`
	City        string `json:"city"`
	Country     string `json:"country"`
}

func (m *UpsertAddressModel) Validate() error {
	if m.Street == "" {
		return InvalidStreet
	}
	if m.HouseNumber == "" {
		return InvalidHouseNumber
	}
	if m.PostalCode == "" {
		return InvalidPostalCode
	}
	if m.City == "" {
		return InvalidCity
	}
	if m.Country == "" {
		return InvalidCountry
	}
	return nil
}

type AddressModel struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Street      string    `json:"street"`
	HouseNumber string    `json:"houseNumber"`
	PostalCode  string    `json:"postalCode"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
}
