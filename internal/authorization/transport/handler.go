package transport

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/transport/models"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

func (h *Handler) CreateRule(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	actorID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	var body models.CreateRuleRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	viewerID, err := uuid.FromString(body.ViewerID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	rule, err := h.service.CreateRule(r.Context(), actorID, actorID, viewerID, body.Resource, body.Effect)
	if err != nil {
		status := mapError(err)
		response.WriteError(w, status, err)
		return
	}

	response.Write(w, http.StatusCreated, models.ToRuleResponse(rule))
}

func (h *Handler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	actorID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	ruleID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	if err := h.service.DeleteRule(r.Context(), actorID, ruleID); err != nil {
		status := mapError(err)
		response.WriteError(w, status, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListRules(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	ownerID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	rules, err := h.service.ListRules(r.Context(), ownerID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToRuleResponses(rules))
}

func (h *Handler) GetRulesByViewer(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	ownerID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	viewerID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	rules, err := h.service.ListRulesByViewer(r.Context(), ownerID, viewerID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToRuleResponses(rules))
}

func mapError(err error) int {
	switch {
	case errors.Is(err, domain.NotOwner):
		return http.StatusForbidden
	case errors.Is(err, domain.SelfRule):
		return http.StatusBadRequest
	case errors.Is(err, domain.InvalidResource):
		return http.StatusUnprocessableEntity
	case errors.Is(err, domain.InvalidEffect):
		return http.StatusUnprocessableEntity
	case errors.Is(err, domain.RuleNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func (h *Handler) PatchResourcePrivacy(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	actorID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	resource := chi.URLParam(r, "resource")
	if !domain.ValidResource(resource) {
		response.WriteError(w, http.StatusUnprocessableEntity, domain.InvalidResource)
		return
	}

	var body models.PatchResourcePrivacyRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.service.SetResourcePrivacy(r.Context(), actorID, resource, body.IsPrivate); err != nil {
		status := mapError(err)
		response.WriteError(w, status, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetResourcePrivacySettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	actorID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	settingsMap, err := h.service.GetResourcePrivacySettings(r.Context(), actorID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var responses []models.ResourcePrivacyResponse
	for res, isPrivate := range settingsMap {
		responses = append(responses, models.ResourcePrivacyResponse{
			Resource:  res,
			IsPrivate: isPrivate,
		})
	}

	// Ensure we don't return null but empty array if no settings
	if responses == nil {
		responses = []models.ResourcePrivacyResponse{}
	}

	response.Write(w, http.StatusOK, responses)
}
