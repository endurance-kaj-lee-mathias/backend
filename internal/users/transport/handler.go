package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/transport/models"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.AddUserModel

	if err := request.Decode(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	usr, err := h.service.AddUser(r.Context(), req.Email, req.Roles)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response.Write(w, http.StatusCreated, models.ToModel(usr))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := domain.ParseId(chi.URLParam(r, "id"))

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, InvalidId)
		return
	}

	usr, err := h.service.GetByID(r.Context(), id)

	if err != nil {
		response.WriteError(w, http.StatusNotFound, NotFound)
		return
	}

	response.Write(w, http.StatusOK, models.ToModel(usr))
}
