package transport

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
)

type ConversationLister interface {
	ListConversationIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type Handler struct {
	manager        *application.Manager
	conversations  ConversationLister
	authenticate   func(*http.Request) (*auth.Claims, error)
	allowedOrigins []string
}

func NewHandler(manager *application.Manager, conversations ConversationLister, authenticate func(*http.Request) (*auth.Claims, error), allowedOrigins []string) *Handler {
	return &Handler{
		manager:        manager,
		conversations:  conversations,
		authenticate:   authenticate,
		allowedOrigins: allowedOrigins,
	}
}
