package ws

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/transport"
)

func Wire(idp config.Idp, allowedOrigins []string) *transport.Handler {
	manager := application.NewManager()
	return transport.NewHandler(manager, auth.ValidateToken(idp), allowedOrigins)
}
