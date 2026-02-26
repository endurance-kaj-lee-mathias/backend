package transport

import (
	"net/http"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
)

type Handler struct {
	manager        *application.Manager
	authenticate   func(*http.Request) (*auth.Claims, error)
	allowedOrigins []string
}

func NewHandler(manager *application.Manager, authenticate func(*http.Request) (*auth.Claims, error), allowedOrigins []string) *Handler {
	return &Handler{
		manager:        manager,
		authenticate:   authenticate,
		allowedOrigins: allowedOrigins,
	}
}
