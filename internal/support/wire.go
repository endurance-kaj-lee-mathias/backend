package support

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/transport"
)

func Wire(db *sql.DB) *transport.Handler {
	repo := infrastructure.NewRepository(db)
	service := application.NewService(repo)
	return transport.NewHandler(service)
}
