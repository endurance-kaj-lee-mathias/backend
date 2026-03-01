package mood

import (
	"context"
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/transport"
)

type Notifier interface {
	Notify(ctx context.Context, deviceToken string) error
}

func Wire(db *sql.DB, enc encryption.Service, notifier Notifier) (*transport.Handler, *application.Scheduler) {
	repo := infrastructure.NewRepository(db)
	userKeyReader := infrastructure.NewUserKeyReader(db, enc)
	service := application.NewService(repo, userKeyReader, enc)
	scheduler := application.NewScheduler(repo, notifier, enc.Hash("VETERAN"))
	handler := transport.NewHandler(service)
	return handler, scheduler
}
