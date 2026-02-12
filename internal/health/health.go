package health

import (
	"database/sql"
	"net/http"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

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

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	dbStatus := "available"
	if err := h.db.PingContext(r.Context()); err != nil {
		dbStatus = "unavailable"
	}

	status := Status{
		Backend:  "UP",
		Database: dbStatus,
	}

	statusCode := http.StatusOK
	if dbStatus == "unavailable" {
		statusCode = http.StatusServiceUnavailable
	}

	response.Write(w, statusCode, status)
}
