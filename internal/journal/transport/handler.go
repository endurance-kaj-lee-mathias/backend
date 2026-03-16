package transport

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/pagination"
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

	limit, offset := pagination.ParsePagination(r)

	report, err := h.service.GetJournal(r.Context(), viewerID, veteranID, limit, offset)
	if err != nil {
		if errors.Is(err, infrastructure.UserNotFound) {
			response.WriteError(w, http.StatusNotFound, VeteranNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToJournalResponse(report, limit, offset))
}
