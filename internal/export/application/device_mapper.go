package application

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) mapDevices(ents []entities.DeviceExportEntity) []domain.DeviceData {
	result := make([]domain.DeviceData, 0, len(ents))

	for _, ent := range ents {
		result = append(result, domain.DeviceData{
			Token:     ent.DeviceToken,
			Platform:  ent.Platform,
			CreatedAt: ent.CreatedAt,
		})
	}

	return result
}
