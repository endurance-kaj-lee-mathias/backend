package calendar

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/transport"
)

func Wire(db *sql.DB, minUrgentMinutes int) *transport.Handler {
	repo := infrastructure.NewRepository(db)
	service := application.NewService(repo, minUrgentMinutes)
	return transport.NewHandler(service)
}
