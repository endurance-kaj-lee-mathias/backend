package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/transport/models"
)

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	veteranId, err := domain.ParseVeteranId(chi.URLParam(r, "id"))

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, InvalidId)
		return
	}

	var req models.AddSupportModel

	if err := request.Decode(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	supportID, err := domain.ParseMemberId(req.ID)

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
	id, err := domain.ParseVeteranId(chi.URLParam(r, "id"))

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, InvalidId)
		return
	}

	mems, err := h.service.GetAll(r.Context(), id)

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToModels(mems))
}
