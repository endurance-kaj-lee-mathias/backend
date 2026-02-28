package mood

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/transport"
)

func Wire(db *sql.DB, enc encryption.Service) (*transport.Handler, *application.Scheduler) {
	repo := infrastructure.NewRepository(db)
	userKeyReader := infrastructure.NewUserKeyReader(db, enc)
	phoneClient := config.NewPhoneClient(config.PhoneServiceURL())
	service := application.NewService(repo, userKeyReader, enc)
	scheduler := application.NewScheduler(repo, phoneClient, enc.Hash("VETERAN"))
	handler := transport.NewHandler(service)
	return handler, scheduler
}
