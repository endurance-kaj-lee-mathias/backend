package main

import (
	"database/sql"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
)

type server struct {
	config config.Config
	db     *sql.DB
}
