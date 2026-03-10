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
	case errors.Is(err, domain.SelfInvite):
		return http.StatusBadRequest, err
	case errors.Is(err, domain.DuplicatePendingInvite),
		errors.Is(err, domain.AlreadyAccepted):
		return http.StatusConflict, err
	case errors.Is(err, domain.NotReceiver):
		return http.StatusForbidden, err
	case errors.Is(err, domain.InviteNotFound),
		errors.Is(err, infrastructure.InviteNotFound):
		return http.StatusNotFound, err
	case errors.Is(err, infrastructure.UserNotFound):
		return http.StatusNotFound, err
	default:
		return http.StatusInternalServerError, err
	}
}
