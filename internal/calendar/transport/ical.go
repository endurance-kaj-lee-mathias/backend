package transport

import (
	ics "github.com/arran4/golang-ical"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

func buildCalendar(events []domain.Event) *ics.Calendar {
	cal := ics.NewCalendar()
	cal.SetVersion("2.0")
	cal.SetProductId("-//Endurance//Calendar//EN")

	for _, e := range events {
		vevent := cal.AddEvent("event-" + e.ID + "@endurance")
		vevent.SetSummary(e.Title)
		vevent.SetStartAt(e.StartTime)
		vevent.SetEndAt(e.EndTime)
		vevent.SetDtStampTime(e.UpdatedAt)
	}

	return cal
}
