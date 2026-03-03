package stress

import (
	"database/sql"
	"net/http"
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/transport"
)

func Wire(db *sql.DB, enc encryption.Service, algoBaseURL string, algoAPIKey string) *transport.Handler {
	repo := infrastructure.NewRepository(db)
	userKeyReader := infrastructure.NewUserKeyReader(db, enc)
	algoClient := infrastructure.NewAlgoClient(algoBaseURL, algoAPIKey, &http.Client{Timeout: 3 * time.Second})
	service := application.NewService(repo, userKeyReader, algoClient, enc)
	return transport.NewHandler(service)
}
