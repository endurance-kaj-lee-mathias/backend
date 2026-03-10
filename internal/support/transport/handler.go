package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/transport/models"
)

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

func (h *Handler) PostInvite(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	senderID, err := domain.ParseMemberId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	var body models.CreateInviteRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	inv, err := h.service.SendInvite(r.Context(), senderID, body.Username, body.Note)
	if err != nil {
		status, errMsg := mapInviteError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	response.Write(w, http.StatusCreated, models.ToInviteModel(inv))
}

func (h *Handler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	callerID, err := domain.ParseMemberId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	inviteID, err := domain.ParseInviteId(chi.URLParam(r, "inviteId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	inv, err := h.service.AcceptInvite(r.Context(), callerID, inviteID)
	if err != nil {
		status, errMsg := mapInviteError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	response.Write(w, http.StatusOK, models.ToInviteModel(inv))
}

func (h *Handler) DeclineInvite(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	callerID, err := domain.ParseMemberId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	inviteID, err := domain.ParseInviteId(chi.URLParam(r, "inviteId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	if err := h.service.DeclineInvite(r.Context(), callerID, inviteID); err != nil {
		status, errMsg := mapInviteError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListInvites(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	callerID, err := domain.ParseMemberId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	incoming, outgoing, err := h.service.ListInvites(r.Context(), callerID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.InviteListModel{
		Incoming: models.ToInviteModels(incoming),
		Outgoing: models.ToInviteModels(outgoing),
	})
}

func (h *Handler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	callerID, err := domain.ParseMemberId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	friendIdStr := chi.URLParam(r, "friendId")
	friendId, err := domain.ParseMemberId(friendIdStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	if err := h.service.DeleteFriend(r.Context(), callerID, friendId); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
