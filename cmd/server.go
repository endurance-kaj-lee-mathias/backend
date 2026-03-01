package main

import (
	"context"
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
)

type notifier interface {
	Notify(ctx context.Context, deviceToken string) error
}

type server struct {
	config   config.Config
	idp      config.Idp
	db       *sql.DB
	enc      encryption.Service
	notifier notifier
}
