package calendar

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/transport"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
)

func Wire(db *sql.DB, enc encryption.Service, minUrgentMinutes int, supportService application.SupportService) *transport.Handler {
	repo := infrastructure.NewRepository(db)
	service := application.NewService(repo, enc, minUrgentMinutes, supportService)
	return transport.NewHandler(service)
}
