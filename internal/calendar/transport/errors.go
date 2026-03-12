package transport

import (
	"errors"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

var Unauthorized = errors.New("unauthorized")
var InvalidId = errors.New("id is invalid")
var CalendarGenerationFailed = errors.New("failed to generate calendar")

func mapError(err error) (int, error) {
	switch {
	case errors.Is(err, domain.OnlyProviderCanManageSlots),
		errors.Is(err, domain.OnlyVeteranCanBook),
		errors.Is(err, domain.NotSlotOwner),
		errors.Is(err, domain.NotAppointmentParticipant):
		return 403, err
	case errors.Is(err, domain.SlotNotFound),
		errors.Is(err, domain.AppointmentNotFound):
		return 404, err
	case errors.Is(err, domain.SlotAlreadyBooked),
		errors.Is(err, domain.SlotOverlap),
		errors.Is(err, domain.CannotDeleteBookedSlot):
		return 409, err
	case errors.Is(err, domain.InvalidTimeRange),
		errors.Is(err, domain.SlotInPast),
		errors.Is(err, domain.NormalCannotBookUrgent),
		errors.Is(err, domain.UrgentRequiresUrgentSlot),
		errors.Is(err, domain.SlotInPastCannotBook),
		errors.Is(err, domain.InsufficientUrgentSlots):
		return 422, err
	default:
		return 500, err
	}
}
