package transport

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

func (h *Handler) ExportUserData(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	export, err := h.service.ExportUserData(r.Context(), userID)
	if err != nil {
		slog.Error("failed to export user data", "userId", userID.String(), "error", err)
		response.WriteError(w, http.StatusInternalServerError, ExportFailed)
		return
	}

	jsonBytes, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		slog.Error("failed to marshal export data", "userId", userID.String(), "error", err)
		response.WriteError(w, http.StatusInternalServerError, ExportFailed)
		return
	}

	filename := fmt.Sprintf("gdpr_export_%s_%s.json", userID.String(), time.Now().UTC().Format("20060102_150405"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
