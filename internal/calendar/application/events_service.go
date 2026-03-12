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
	providerCache := make(map[uuid.UUID]string)

	for _, ent := range ents {
		title, ok := providerCache[ent.ID]

		if !ok {
			providerKey, err := s.enc.DecryptUserKey(ent.ProviderEncryptedUserKey)
			if err != nil {
				return nil, err
			}

			firstNameBytes, err := s.enc.Decrypt(ent.ProviderEncryptedFirstName, providerKey)
			if err != nil {
				return nil, err
			}

			lastNameBytes, err := s.enc.Decrypt(ent.ProviderEncryptedLastName, providerKey)
			if err != nil {
				return nil, err
			}

			title = fmt.Sprintf("Appointment with %s %s", string(firstNameBytes), string(lastNameBytes))
			providerCache[ent.ID] = title
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
