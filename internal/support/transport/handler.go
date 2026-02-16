package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/transport/models"
)

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	veteranIdStr := chi.URLParam(r, "id")

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
		status, errMsg := mapAddMemberError(err)
		response.WriteError(w, status, errMsg)
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

func (h *Handler) DeleteSupporter(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	veteranId, err := domain.ParseVeteranId(claims.Sub)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	supportIdStr := chi.URLParam(r, "supportId")
	supportId, err := domain.ParseMemberId(supportIdStr)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	if err := h.service.DeleteSupporter(r.Context(), veteranId, supportId); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
