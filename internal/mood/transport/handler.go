package transport

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
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

	entries, err := h.service.GetEntriesByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, infrastructure.UserNotFound) {
			response.WriteError(w, http.StatusNotFound, UserNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToResponseList(entries))
}

func (h *Handler) GetEntriesByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := domain.NewUserId(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	entries, err := h.service.GetEntriesByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, infrastructure.UserNotFound) {
			response.WriteError(w, http.StatusNotFound, UserNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToResponseList(entries))
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
