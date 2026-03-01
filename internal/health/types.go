package health

import (
	"database/sql"

	"firebase.google.com/go/v4/messaging"
)

type Handler struct {
	db       *sql.DB
	firebase *messaging.Client
}

func NewHandler(db *sql.DB, firebase *messaging.Client) *Handler {
	return &Handler{db: db, firebase: firebase}
}

type Status struct {
	Backend  string `json:"backend"`
	Database string `json:"database"`
	Firebase string `json:"firebase"`
}
