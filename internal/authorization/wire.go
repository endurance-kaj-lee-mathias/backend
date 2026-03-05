package authorization

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/transport"
)

func Wire(db *sql.DB) (*transport.Handler, application.Service) {
	repo := infrastructure.NewRepository(db)
	service := application.NewService(repo)
	handler := transport.NewHandler(service)
	return handler, service
}
