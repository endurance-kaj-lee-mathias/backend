package chats

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/transport"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	wsapplication "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
)

func Wire(db *sql.DB, enc encryption.Service, notifier application.Notifier, broadcaster wsapplication.Broadcaster) (*transport.Handler, application.Service) {
	repo := infrastructure.NewRepository(db, enc)
	svc := application.NewService(repo, enc, notifier, broadcaster)
	return transport.NewHandler(svc), svc
}
