package health

import (
	"net/http"

	"firebase.google.com/go/v4/messaging"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	dbStatus := "available"
	if err := h.db.PingContext(r.Context()); err != nil {
		dbStatus = "unavailable"
	}

	firebaseStatus := "available"
	_, err := h.firebase.SendEachDryRun(r.Context(), []*messaging.Message{
		{Topic: "health-check"},
	})
	if err != nil {
		firebaseStatus = "unavailable"
	}

	status := Status{
		Backend:  "UP",
		Database: dbStatus,
		Firebase: firebaseStatus,
	}

	statusCode := http.StatusOK
	if dbStatus == "unavailable" || firebaseStatus == "unavailable" {
		statusCode = http.StatusServiceUnavailable
	}

	response.Write(w, statusCode, status)
}
