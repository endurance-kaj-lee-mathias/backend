package domain

import "time"

type Address struct {
	ID         AddressId `json:"id"`
	UserID     UserId    `json:"userId"`
	Street     string    `json:"street"`
	Locality   string    `json:"locality"`
	Region     string    `json:"region"`
	PostalCode string    `json:"postalCode"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"createdAt"`
}

func NewAddress(id AddressId, userID UserId, street string, locality string, region string, postalCode string, country string) Address {
	return Address{
		ID:         id,
		UserID:     userID,
		Street:     street,
		Locality:   locality,
		Region:     region,
		PostalCode: postalCode,
		Country:    country,
		CreatedAt:  time.Now().UTC(),
	}
}
