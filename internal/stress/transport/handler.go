package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/transport/models"
)

const defaultLimit = 20
const defaultOffset = 0

func parsePagination(r *http.Request) (limit, offset int) {
	limit = defaultLimit
	offset = defaultOffset

	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}

func (h *Handler) IngestSample(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	var body models.StressSampleRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if body.UserID != claims.Sub {
		response.WriteError(w, http.StatusForbidden, UserIdMismatch)
		return
	}

	userID, err := uuid.FromString(body.UserID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	sample, err := domain.NewStressSample(
		userID,
		body.TimestampUTC,
		body.WindowMinutes,
		body.MeanHR,
		body.RMSSDms,
		body.RestingHR,
		body.Steps,
		body.SleepDebtHours,
	)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.service.IngestSample(r.Context(), sample); err != nil {
		if errors.Is(err, infrastructure.UserNotFound) {
			response.WriteError(w, http.StatusNotFound, UserNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetLatestSampleTimestamp(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	userID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	ts, err := h.service.GetLatestSampleTimestamp(r.Context(), userID)
	if err != nil {
		if errors.Is(err, infrastructure.SampleNotFound) {
			response.WriteError(w, http.StatusNotFound, SampleNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.LatestSampleResponse{Timestamp: ts})
}

func (h *Handler) GetLatestScore(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	userID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	limit, offset := parsePagination(r)

	scores, total, err := h.service.GetScoresPaginated(r.Context(), userID, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, response.NewPaginated(models.ToStressScoreResponseList(scores), total, limit, offset))
}

func (h *Handler) DeleteMySamples(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	userID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, InvalidId)
		return
	}

	if err := h.service.DeleteMySamples(r.Context(), userID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetLatestScoreByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	limit, offset := parsePagination(r)

	scores, total, err := h.service.GetScoresPaginated(r.Context(), userID, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, response.NewPaginated(models.ToStressScoreResponseList(scores), total, limit, offset))
}
