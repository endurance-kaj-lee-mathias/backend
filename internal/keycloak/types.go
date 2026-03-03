package keycloak

import (
	"context"
	"sync"
	"time"
)

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

type client struct {
	baseURL       string
	realm         string
	adminUser     string
	adminPassword string

	mu          sync.Mutex
	token       string
	tokenExpiry time.Time
}

func NewClient(baseURL string, realm string, adminUser string, adminPassword string) Client {
	return &client{
		baseURL:       baseURL,
		realm:         realm,
		adminUser:     adminUser,
		adminPassword: adminPassword,
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type keycloakUser struct {
	FirstName  string              `json:"firstName"`
	LastName   string              `json:"lastName"`
	Email      string              `json:"email"`
	Username   string              `json:"username"`
	Attributes map[string][]string `json:"attributes"`
}
