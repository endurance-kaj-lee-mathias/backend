package support

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/transport"
)

func Wire(db *sql.DB, enc encryption.Service, authz application.AuthzRevoker, notifier application.Notifier) *transport.Handler {
	repo := infrastructure.NewRepository(db, enc)
	inviteRepo := infrastructure.NewInviteRepository(db, enc)
	userRoleRead := infrastructure.NewUserRoleReader(db, enc)
	service := application.NewService(repo, inviteRepo, userRoleRead, enc, authz, notifier)
	return transport.NewHandler(service)
}
