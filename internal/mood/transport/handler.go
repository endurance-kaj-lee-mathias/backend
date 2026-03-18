package transport

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/pagination"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/transport/models"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

func (h *Handler) UpsertMoodEntry(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	var body models.MoodEntryRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	userID, err := domain.NewUserId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, InvalidId)
		return
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, InvalidDate)
		return
	}

	entry, err := domain.NewMoodEntry(userID, date, body.MoodScore, body.Notes)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.service.UpsertMoodEntry(r.Context(), entry); err != nil {
		if errors.Is(err, infrastructure.UserNotFound) {
			response.WriteError(w, http.StatusNotFound, UserNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetMyEntries(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	userID, err := domain.NewUserId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, InvalidId)
		return
	}

	limit, offset := pagination.ParsePagination(r)

	entries, total, err := h.service.GetEntriesByUserID(r.Context(), userID, offset)
	if err != nil {
		if errors.Is(err, infrastructure.UserNotFound) {
			response.WriteError(w, http.StatusNotFound, UserNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, response.NewPaginated(models.ToResponseList(entries), total, limit, offset))
}

func (h *Handler) GetVeteransSupport(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	memberID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, InvalidId)
		return
	}

	summaries, err := h.service.GetVeteransSupport(r.Context(), memberID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToVeteranSupportResponseList(summaries))
}

func (h *Handler) GetEntriesByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := domain.NewUserId(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	limit, offset := pagination.ParsePagination(r)

	entries, total, err := h.service.GetEntriesByUserID(r.Context(), userID, offset)
	if err != nil {
		if errors.Is(err, infrastructure.UserNotFound) {
			response.WriteError(w, http.StatusNotFound, UserNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, response.NewPaginated(models.ToResponseList(entries), total, limit, offset))
}

func (h *Handler) GetTodayEntry(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	userID, err := domain.NewUserId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, InvalidId)
		return
	}

	entry, err := h.service.GetTodayEntry(r.Context(), userID)
	if err != nil {
		if errors.Is(err, infrastructure.MoodEntryNotFound) {
			response.WriteError(w, http.StatusNotFound, MoodEntryNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToResponse(*entry))
}

func (h *Handler) UpdateMoodEntry(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	entryID, err := domain.MoodIdFromString(chi.URLParam(r, "entryId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidEntryId)
		return
	}

	existing, err := h.service.GetEntryByID(r.Context(), entryID)
	if err != nil {
		if errors.Is(err, infrastructure.MoodEntryNotFound) {
			response.WriteError(w, http.StatusNotFound, MoodEntryNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if existing.UserID.String() != claims.Sub {
		response.WriteError(w, http.StatusForbidden, Forbidden)
		return
	}

	var body models.MoodEntryRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, InvalidDate)
		return
	}

	updated := domain.MoodEntry{
		ID:        existing.ID,
		UserID:    existing.UserID,
		Date:      date,
		MoodScore: body.MoodScore,
		Notes:     body.Notes,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: existing.UpdatedAt,
	}

	if err := h.service.UpdateMoodEntry(r.Context(), updated); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteMoodEntry(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	entryID, err := domain.MoodIdFromString(chi.URLParam(r, "entryId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidEntryId)
		return
	}

	existing, err := h.service.GetEntryByID(r.Context(), entryID)
	if err != nil {
		if errors.Is(err, infrastructure.MoodEntryNotFound) {
			response.WriteError(w, http.StatusNotFound, MoodEntryNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if existing.UserID.String() != claims.Sub {
		response.WriteError(w, http.StatusForbidden, Forbidden)
		return
	}

	if err := h.service.DeleteMoodEntry(r.Context(), entryID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteMyEntries(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	userID, err := domain.NewUserId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, InvalidId)
		return
	}

	if err := h.service.DeleteMyMoodEntries(r.Context(), userID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
