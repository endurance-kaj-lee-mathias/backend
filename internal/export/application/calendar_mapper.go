package application

import (
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) decryptAppointments(ents []entities.AppointmentExportEntity) ([]domain.AppointmentData, error) {
	result := make([]domain.AppointmentData, 0, len(ents))
	providerCache := make(map[uuid.UUID]string)

	for _, ent := range ents {
		providerUsername, ok := providerCache[ent.ProviderID]
		if !ok {
			providerKey, err := s.enc.DecryptUserKey(ent.ProviderEncryptedUserKey)
			if err != nil {
				return nil, err
			}

			username, err := s.enc.Decrypt(ent.ProviderEncryptedUsername, providerKey)
			if err != nil {
				return nil, err
			}

			providerUsername = string(username)
			providerCache[ent.ProviderID] = providerUsername
		}

		var title *string
		if ent.AppointmentTitle.Valid {
			title = &ent.AppointmentTitle.String
		}

		result = append(result, domain.AppointmentData{
			ID:               ent.ID.String(),
			SlotID:           ent.SlotID.String(),
			VeteranID:        ent.VeteranID.String(),
			ProviderID:       ent.ProviderID.String(),
			ProviderUsername: providerUsername,
			Title:            title,
			Status:           ent.Status,
			StartTime:        ent.StartTime,
			EndTime:          ent.EndTime,
			IsUrgent:         ent.IsUrgent,
			CreatedAt:        ent.CreatedAt,
			UpdatedAt:        ent.UpdatedAt,
		})
	}

	return result, nil
}

func (s *service) mapSlots(ents []entities.SlotExportEntity) []domain.SlotData {
	result := make([]domain.SlotData, 0, len(ents))

	for _, ent := range ents {
		result = append(result, domain.SlotData{
			ID:        ent.ID.String(),
			StartTime: ent.StartTime,
			EndTime:   ent.EndTime,
			IsUrgent:  ent.IsUrgent,
			IsBooked:  ent.IsBooked,
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
		})
	}

	return result
}
