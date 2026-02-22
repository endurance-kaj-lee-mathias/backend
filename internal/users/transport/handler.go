package transport

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/request"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/transport/models"
)

func (h *Handler) PatchPhoneNumber(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	id, err := domain.ParseId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	var body models.PhoneNumberModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdatePhoneNumber(r.Context(), id, body.PhoneNumber); err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpsertAddress(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	id, err := domain.ParseId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	var body models.UpsertAddressModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	addr, err := h.service.UpsertAddress(r.Context(), id, body.Street, body.HouseNumber, body.PostalCode, body.City, body.Country)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.AddressModel{
		ID:          addr.ID.UUID,
		UserID:      addr.UserID.UUID,
		Street:      addr.Street,
		HouseNumber: addr.HouseNumber,
		PostalCode:  addr.PostalCode,
		City:        addr.City,
		Country:     addr.Country,
	})
}

func (h *Handler) GetAddress(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	id, err := domain.ParseId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	addr, err := h.service.GetAddress(r.Context(), id)
	if err != nil {
		if errors.Is(err, infrastructure.AddressNotFound) {
			response.WriteError(w, http.StatusNotFound, AddressNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.AddressModel{
		ID:          addr.ID.UUID,
		UserID:      addr.UserID.UUID,
		Street:      addr.Street,
		HouseNumber: addr.HouseNumber,
		PostalCode:  addr.PostalCode,
		City:        addr.City,
		Country:     addr.Country,
	})
}

func (h *Handler) GetOrCreate(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())

	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	id, err := domain.ParseId(claims.Sub)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	roles := make([]domain.Role, 0, len(claims.Roles))
	for _, role := range claims.Roles {
		roles = append(roles, domain.Role(role))
	}

	usr, err := h.service.GetOrCreate(
		r.Context(),
		id,
		claims.Email,
		claims.FirstName,
		claims.LastName,
		roles,
	)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToModel(usr))
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

func (h *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserClaims(r.Context())

	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	id, err := domain.ParseId(claims.Sub)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
