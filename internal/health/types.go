package health

import "database/sql"

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db}
}

type Status struct {
	Backend  string `json:"backend"`
	Database string `json:"database"`
}
