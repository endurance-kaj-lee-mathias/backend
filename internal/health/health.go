package health

import (
	"net/http"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

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
