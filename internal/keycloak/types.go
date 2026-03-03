package keycloak

import "context"

type Client interface {
	UpdateUser(ctx context.Context, userID string, update UserUpdate) error
}

type UserUpdate struct {
	FirstName   string
	LastName    string
	Email       string
	Username    string
	PhoneNumber *string
	Street      *string
	Locality    *string
	Region      *string
	PostalCode  *string
	Country     *string
}
