package transport

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/transport/models"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"net/http"
)

func (h *Handler) GetMessage(w http.ResponseWriter, r *http.Request) {
	msg, err := h.service.GetMessage(r.Context())

	if err != nil {
		response.WriteError(w, http.StatusNotFound, NotFound)
		return
	}

	model := models.ToModel(msg)
	response.Write(w, http.StatusOK, model)
}
