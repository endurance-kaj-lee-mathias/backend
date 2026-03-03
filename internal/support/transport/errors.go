package transport

import (
	"errors"
	"net/http"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
)

var InvalidId = errors.New("id is invalid")
var Unauthorized = errors.New("unauthorized")

func mapInviteError(err error) (int, error) {
	switch {
	case errors.Is(err, domain.ErrSelfInvite):
		return http.StatusBadRequest, err
	case errors.Is(err, domain.ErrDuplicatePendingInvite),
		errors.Is(err, domain.ErrAlreadyAccepted):
		return http.StatusConflict, err
	case errors.Is(err, domain.ErrNotReceiver):
		return http.StatusForbidden, err
	case errors.Is(err, domain.ErrInviteNotFound),
		errors.Is(err, infrastructure.InviteNotFound):
		return http.StatusNotFound, err
	case errors.Is(err, domain.VeteranMustHaveVeteranRole),
		errors.Is(err, domain.SupporterMustBeAbleToSupport):
		return http.StatusForbidden, err
	case errors.Is(err, infrastructure.UserNotFound):
		return http.StatusNotFound, err
	default:
		return http.StatusInternalServerError, err
	}
}
