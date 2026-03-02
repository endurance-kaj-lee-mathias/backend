package domain

import (
	"strings"
	"time"

	userdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

func ValidateSlotTimes(start, end time.Time) error {
	if !end.After(start) {
		return InvalidTimeRange
	}
	if start.Before(time.Now().UTC()) {
		return SlotInPast
	}
	return nil
}

func ValidateBooking(slot Slot, urgent bool) error {
	if slot.IsBooked {
		return SlotAlreadyBooked
	}
	if slot.StartTime.Before(time.Now().UTC()) {
		return SlotInPastCannotBook
	}
	if !urgent && slot.IsUrgent {
		return NormalCannotBookUrgent
	}
	if urgent && !slot.IsUrgent {
		return UrgentRequiresUrgentSlot
	}
	return nil
}

func HasProviderRole(roles []string) bool {
	for _, r := range roles {
		upper := strings.ToUpper(r)
		if upper == string(userdomain.RoleTherapist) || upper == string(userdomain.RoleSupport) {
			return true
		}
	}
	return false
}

func HasVeteranRole(roles []string) bool {
	for _, r := range roles {
		if strings.ToUpper(r) == string(userdomain.RoleVeteran) {
			return true
		}
	}
	return false
}
