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
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpdatePhoneNumberModel
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
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
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

	addr, err := h.service.UpsertAddress(r.Context(), id, body.Street, body.Locality, body.Region, body.PostalCode, body.Country)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.Write(w, http.StatusOK, models.ToAddressModel(addr))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	id, claims, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	roleName := ""
	for _, aud := range claims.Audiences {
		if aud == h.mobileClientID {
			roleName = "veteran"
			break
		}
		if aud == h.webClientID {
			roleName = "support-group"
			break
		}
	}

	if roleName == "" {
		response.WriteError(w, http.StatusForbidden, ClientNotEligible)
		return
	}

	if err := h.service.AssignRole(r.Context(), id, roleName); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetOrCreate(w http.ResponseWriter, r *http.Request) {
	id, claims, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	roles := make([]domain.Role, 0, len(claims.Roles))
	for _, role := range claims.Roles {
		roles = append(roles, domain.Role(role))
	}

	usr, err := h.service.GetOrCreate(r.Context(), id, claims.Email, claims.Username, claims.FirstName, claims.LastName, claims.PhoneNumber, claims.Address.StreetAddress, claims.Address.Locality, claims.Address.Region, claims.Address.PostalCode, claims.Address.Country, roles)
	if err != nil {
		if errors.Is(err, infrastructure.UsernameAlreadyExists) {
			response.WriteError(w, http.StatusConflict, UsernameAlreadyExists)
			return
		}
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	addr, ok := h.optionalAddress(r.Context(), w, id)
	if !ok {
		return
	}

	response.Write(w, http.StatusOK, models.ToModel(usr, addr))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := domain.ParseId(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return
	}

	usr, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, NotFound)
		return
	}

	addr, ok := h.optionalAddress(r.Context(), w, id)
	if !ok {
		return
	}

	response.Write(w, http.StatusOK, models.ToModel(usr, addr))
}

func (h *Handler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	targetID, ok := auth.GetTargetID(r.Context())

	var usr domain.User
	var err error

	if ok {
		usr, err = h.service.GetByID(r.Context(), domain.UserId{UUID: targetID})
	} else {
		username := chi.URLParam(r, "username")
		usr, err = h.service.GetByUsername(r.Context(), username)
	}

	if err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}

		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	addr, ok := h.optionalAddress(r.Context(), w, usr.ID)
	if !ok {
		return
	}

	response.Write(w, http.StatusOK, models.ToModel(usr, addr))
}

func (h *Handler) PatchFirstName(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpdateFirstNameModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateFirstName(r.Context(), id, body.FirstName); err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchLastName(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpdateLastNameModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateLastName(r.Context(), id, body.LastName); err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchIntroduction(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpdateIntroductionModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateIntroduction(r.Context(), id, body.Introduction); err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchAbout(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpdateAboutModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateAbout(r.Context(), id, body.About); err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchImage(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpdateImageModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateImage(r.Context(), id, body.Image); err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchPrivacy(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpdatePrivacyModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdatePrivacy(r.Context(), id, body.IsPrivate); err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PutDevice(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpsertDeviceModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.service.UpsertDevice(r.Context(), id, body.Token, body.Platform); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	_, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.DeleteDeviceModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.service.DeleteDevice(r.Context(), body.Token); err != nil {
		if errors.Is(err, infrastructure.DeviceNotFound) {
			response.WriteError(w, http.StatusNotFound, DeviceNotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PutRiskLevel(w http.ResponseWriter, r *http.Request) {
	id, _, ok := h.authenticatedID(w, r)
	if !ok {
		return
	}

	var body models.UpdateRiskLevelModel
	if err := request.Decode(r, &body); err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.service.UpdateRiskLevel(r.Context(), id, domain.RiskLevel(body.RiskLevel)); err != nil {
		if errors.Is(err, infrastructure.NotFound) {
			response.WriteError(w, http.StatusNotFound, NotFound)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
