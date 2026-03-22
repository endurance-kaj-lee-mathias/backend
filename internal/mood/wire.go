package mood

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/transport"
)

type Notifier interface {
	Notify(ctx context.Context, userID uuid.UUID, deviceTokens []string) error
}

func Wire(db *sql.DB, enc encryption.Service, notifier Notifier, authz application.AuthorizationChecker) (*transport.Handler, *application.Scheduler) {
	repo := infrastructure.NewRepository(db)
	userKeyReader := infrastructure.NewUserKeyReader(db, enc)
	veteranReader := infrastructure.NewVeteranReader(db, enc)
	service := application.NewService(repo, userKeyReader, enc, veteranReader, authz)
	scheduler := application.NewScheduler(repo, notifier, enc.Hash("VETERAN"))
	handler := transport.NewHandler(service)
	return handler, scheduler
}
