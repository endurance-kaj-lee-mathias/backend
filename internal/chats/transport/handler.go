package transport

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/pagination"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/transport/models"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"

	"github.com/gofrs/uuid"
)

func (h *Handler) StartConversation(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	callerID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	var body models.StartConversationRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidRequestBody)
		return
	}

	if body.ParticipantID == uuid.Nil {
		response.WriteError(w, http.StatusBadRequest, models.InvalidParticipantID)
		return
	}

	conv, err := h.service.GetOrCreateConversation(r.Context(), callerID, body.ParticipantID)
	if err != nil {
		if errors.Is(err, application.NoSupportRelationship) {
			response.WriteError(w, http.StatusForbidden, Forbidden)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToConversationModel(conv))
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	senderID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	convID, err := domain.ParseConversationId(chi.URLParam(r, "conversationId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidConversationID)
		return
	}

	var body models.SendMessageRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidRequestBody)
		return
	}

	if body.Content == "" {
		response.WriteError(w, http.StatusBadRequest, models.EmptyContent)
		return
	}

	msg, err := h.service.SendMessage(r.Context(), convID.UUID, senderID, body.Content)
	if err != nil {
		if errors.Is(err, infrastructure.ParticipantNotFound) {
			response.WriteError(w, http.StatusForbidden, Forbidden)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusCreated, models.ToCreatedMessageModel(msg))
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	callerID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	convID, err := domain.ParseConversationId(chi.URLParam(r, "conversationId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidConversationID)
		return
	}

	limit, offset := pagination.ParsePagination(r)

	msgs, err := h.service.GetMessages(r.Context(), convID.UUID, callerID, limit, offset)
	if err != nil {
		if errors.Is(err, infrastructure.ParticipantNotFound) {
			response.WriteError(w, http.StatusForbidden, Forbidden)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToMessageModels(msgs))
}

func (h *Handler) GetAllChats(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	callerID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	summaries, err := h.service.GetAllChats(r.Context(), callerID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToChatSummaryModels(summaries))
}
