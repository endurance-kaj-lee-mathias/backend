package domain

import (
	"time"
)

type ProfileData struct {
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	PhoneNumber  *string   `json:"phoneNumber,omitempty"`
	Roles        []string  `json:"roles"`
	About        string    `json:"about"`
	Introduction string    `json:"introduction"`
	Image        string    `json:"image"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AddressData struct {
	Street     string    `json:"street"`
	Locality   string    `json:"locality"`
	Region     string    `json:"region"`
	PostalCode string    `json:"postalCode"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"createdAt"`
}
