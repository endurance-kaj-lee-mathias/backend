package transport

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/transport/models"
)

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	veteranIdStr := chi.URLParam(r, "id")
	if strings.EqualFold(veteranIdStr, "me") {
		// handle "me" if needed, but for now we follow the instruction that
		// the account that posts it is the support member.
	}

	veteranId, err := domain.ParseVeteranId(veteranIdStr)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	supportID, err := domain.ParseMemberId(claims.Sub)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	mem, err := h.service.AddMember(r.Context(), veteranId, supportID)

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusCreated, models.ToModel(mem))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	id, err := domain.ParseVeteranId(claims.Sub)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	mems, err := h.service.GetAll(r.Context(), id)

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToModels(mems))
}
