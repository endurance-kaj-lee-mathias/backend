package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/transport/models"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

func (h *Handler) GetJournal(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	viewerID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	veteranID, ok := auth.GetTargetID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusBadRequest, InvalidUsername)
		return
	}

	weekOffset := 0

	if v := r.URL.Query().Get("week"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			weekOffset = parsed
		}
	}

	report, err := h.service.GetJournal(r.Context(), viewerID, veteranID, weekOffset)
	if err != nil {
		if errors.Is(err, infrastructure.UserNotFound) {
			response.WriteError(w, http.StatusNotFound, VeteranNotFound)
			return
		}

		if errors.Is(err, domain.NotVeteran) {
			response.WriteError(w, http.StatusForbidden, Forbidden)
			return
		}

		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	journalResponse, err := models.ToJournalResponse(report, weekOffset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, journalResponse)
}
