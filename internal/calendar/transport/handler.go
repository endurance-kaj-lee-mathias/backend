package transport

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/transport/models"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

func (h *Handler) CreateSlot(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	providerID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	var body models.CreateSlotRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	slot, err := h.service.CreateSlot(r.Context(), providerID, claims.Roles, body.StartTime, body.EndTime, body.IsUrgent)
	if err != nil {
		status, errMsg := mapError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	response.Write(w, http.StatusCreated, models.ToSlotModel(slot))
}

func (h *Handler) GetSlots(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, models.InvalidFromParam)
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, models.InvalidToParam)
		return
	}

	var providerID *uuid.UUID
	if pidStr := r.URL.Query().Get("providerId"); pidStr != "" {
		pid, err := uuid.FromString(pidStr)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, models.InvalidProviderIdParam)
			return
		}
		providerID = &pid
	}

	slots, err := h.service.GetSlots(r.Context(), from, to, providerID)
	if err != nil {
		status, errMsg := mapError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	response.Write(w, http.StatusOK, models.ToSlotModels(slots))
}

func (h *Handler) DeleteSlot(w http.ResponseWriter, r *http.Request) {
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

	slotID, err := domain.ParseSlotId(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	if err := h.service.DeleteSlot(r.Context(), userID, claims.Roles, slotID); err != nil {
		status, errMsg := mapError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) BookSlot(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	veteranID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	slotID, err := domain.ParseSlotId(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	var body models.BookSlotRequest
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	appointment, err := h.service.BookSlot(r.Context(), veteranID, claims.Roles, slotID, body.Urgent)
	if err != nil {
		status, errMsg := mapError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	response.Write(w, http.StatusCreated, models.ToAppointmentModel(appointment))
}

func (h *Handler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
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

	appointmentID, err := domain.ParseAppointmentId(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	if err := h.service.CancelAppointment(r.Context(), userID, appointmentID); err != nil {
		status, errMsg := mapError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSlotsByUserID(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, models.InvalidFromParam)
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, models.InvalidToParam)
		return
	}

	providerID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	slots, err := h.service.GetSlots(r.Context(), from, to, &providerID)
	if err != nil {
		status, errMsg := mapError(err)
		response.WriteError(w, status, errMsg)
		return
	}

	response.Write(w, http.StatusOK, models.ToSlotModels(slots))
}

func (h *Handler) ExportCalendar(w http.ResponseWriter, r *http.Request) {
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

	events, err := h.service.GetCalendarEvents(r.Context(), userID)
	if err != nil {
		slog.Error("failed to get calendar events", "userId", userID.String(), "error", err)
		response.WriteError(w, http.StatusInternalServerError, CalendarGenerationFailed)
		return
	}

	cal := buildCalendar(events)

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", `attachment; filename="calendar.ics"`)
	w.WriteHeader(http.StatusOK)

	if err := cal.SerializeTo(w); err != nil {
		slog.Error("failed to serialize calendar", "userId", userID.String(), "error", err)
	}
}

func (h *Handler) FeedCalendar(w http.ResponseWriter, r *http.Request) {
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

	events, err := h.service.GetCalendarEvents(r.Context(), userID)
	if err != nil {
		slog.Error("failed to get calendar feed", "userId", userID.String(), "error", err)
		response.WriteError(w, http.StatusInternalServerError, CalendarGenerationFailed)
		return
	}

	cal := buildCalendar(events)

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", `inline; filename="calendar.ics"`)
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	if err := cal.SerializeTo(w); err != nil {
		slog.Error("failed to serialize calendar feed", "userId", userID.String(), "error", err)
	}
}
