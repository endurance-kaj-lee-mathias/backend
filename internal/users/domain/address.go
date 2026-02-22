package domain

import "time"

type Address struct {
	ID          AddressId `json:"id"`
	UserID      UserId    `json:"userId"`
	Street      string    `json:"street"`
	HouseNumber string    `json:"houseNumber"`
	PostalCode  string    `json:"postalCode"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewAddress(id AddressId, userID UserId, street string, houseNumber string, postalCode string, city string, country string) Address {
	return Address{
		ID:          id,
		UserID:      userID,
		Street:      street,
		HouseNumber: houseNumber,
		PostalCode:  postalCode,
		City:        city,
		Country:     country,
		CreatedAt:   time.Now().UTC(),
	}
}
