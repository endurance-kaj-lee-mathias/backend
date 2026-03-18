package application

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
)

func (s *service) GetCalendarEvents(ctx context.Context, userID uuid.UUID) ([]domain.Event, error) {
	ents, err := s.repo.GetEventsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	events := make([]domain.Event, 0, len(ents))

	for _, ent := range ents {
		var encFirstName, encLastName, encUserKey []byte
		if userID == ent.VeteranID {
			encFirstName = ent.ProviderEncryptedFirstName
			encLastName = ent.ProviderEncryptedLastName
			encUserKey = ent.ProviderEncryptedUserKey
		} else {
			encFirstName = ent.VeteranEncryptedFirstName
			encLastName = ent.VeteranEncryptedLastName
			encUserKey = ent.VeteranEncryptedUserKey
		}

		userKey, err := s.enc.DecryptUserKey(encUserKey)
		if err != nil {
			return nil, err
		}

		firstNameBytes, err := s.enc.Decrypt(encFirstName, userKey)
		if err != nil {
			return nil, err
		}

		lastNameBytes, err := s.enc.Decrypt(encLastName, userKey)
		if err != nil {
			return nil, err
		}

		var title string
		if ent.AppointmentTitle.Valid && ent.AppointmentTitle.String != "" {
			title = ent.AppointmentTitle.String
		} else {
			title = fmt.Sprintf("Appointment with %s %s", string(firstNameBytes), string(lastNameBytes))
		}

		events = append(events, domain.Event{
			ID:        ent.ID.String(),
			Title:     title,
			StartTime: ent.StartTime.UTC(),
			EndTime:   ent.EndTime.UTC(),
			UpdatedAt: ent.UpdatedAt.UTC(),
		})
	}

	return events, nil
}
