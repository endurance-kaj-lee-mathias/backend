package stress

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/transport"
)

func Wire(db *sql.DB, enc encryption.Service) *transport.Handler {
	repo := infrastructure.NewRepository(db)
	userKeyReader := infrastructure.NewUserKeyReader(db, enc)
	service := application.NewService(repo, userKeyReader, enc)
	return transport.NewHandler(service)
}
