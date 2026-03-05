package transport

import (
	"context"
	"errors"
	"net/http"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"

	"github.com/gofrs/uuid"
)

type Handler struct {
	service application.Service
}

func NewHandler(s application.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) authenticatedID(w http.ResponseWriter, r *http.Request) (domain.UserId, *auth.Claims, bool) {
	claims, ok := auth.GetUserClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return domain.UserId{}, nil, false
	}

	id, err := domain.ParseId(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidId)
		return domain.UserId{}, nil, false
	}

	return id, claims, true
}

func (h *Handler) optionalAddress(ctx context.Context, w http.ResponseWriter, id domain.UserId) (*domain.Address, bool) {
	addr, err := h.service.GetAddress(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.AddressNotFound) {
			return nil, true
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return nil, false
	}
	return &addr, true
}

func (h *Handler) ResolveUsername(ctx context.Context, username string) (uuid.UUID, error) {
	usr, err := h.service.GetByUsername(ctx, username)
	if err != nil {
		return uuid.UUID{}, err
	}
	return usr.ID.UUID, nil
}
