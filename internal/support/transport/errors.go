package transport

import (
	"errors"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
)

var InvalidId = errors.New("id is invalid")
var Unauthorized = errors.New("unauthorized")
var NotFound = errors.New("member was not found")

func mapAddMemberError(err error) (int, error) {
	switch {
	case errors.Is(err, domain.VeteranMustHaveVeteranRole),
		errors.Is(err, domain.SupporterMustBeAbleToSupport):
		return 403, err
	case errors.Is(err, domain.SelfSupportNotAllowed):
		return 400, err
	case errors.Is(err, infrastructure.UserNotFound):
		return 404, err
	default:
		return 500, err
	}
}
