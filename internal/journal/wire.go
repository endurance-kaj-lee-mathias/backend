package journal

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/transport"
)

type Handler = transport.Handler

func Wire(db *sql.DB, enc encryption.Service, authz application.AuthorizationChecker) *Handler {
	repo := infrastructure.NewRepository(db)
	svc := application.NewService(repo, enc, authz)
	return transport.NewHandler(svc)
}
