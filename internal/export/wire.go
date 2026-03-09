package export

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/transport"
)

type Handler = transport.Handler

func Wire(db *sql.DB, enc encryption.Service) *Handler {
	repo := infrastructure.NewRepository(db, enc)
	svc := application.NewService(repo, enc)
	return transport.NewHandler(svc)
}
