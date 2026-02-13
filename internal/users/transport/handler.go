package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/transport/models"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.service.CreateUser(r.Context(), req.Email, req.Roles)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	response.Write(w, http.StatusCreated, models.ToModel(*user))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	u, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, NotFound)
		return
	}
	response.Write(w, http.StatusOK, models.ToModel(*u))
}

func (h *Handler) AddSupportMember(w http.ResponseWriter, r *http.Request) {
	veteranIDStr := chi.URLParam(r, "veteranId")
	veteranID, err := uuid.Parse(veteranIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	var req addSupportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	supportID, err := uuid.Parse(req.SupportID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.service.AddSupportMember(r.Context(), veteranID, supportID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListSupportMembers(w http.ResponseWriter, r *http.Request) {
	veteranIDStr := chi.URLParam(r, "veteranId")
	veteranID, err := uuid.Parse(veteranIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	members, err := h.service.ListSupportMembers(r.Context(), veteranID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	out := make([]models.UserModel, 0, len(members))
	for _, m := range members {
		out = append(out, models.ToModel(m))
	}
	response.Write(w, http.StatusOK, out)
}
