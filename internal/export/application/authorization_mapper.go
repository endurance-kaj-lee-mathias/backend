package application

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) mapAuthorizationRules(ents []entities.AuthorizationRuleExportEntity) []domain.AuthorizationRuleData {
	result := make([]domain.AuthorizationRuleData, 0, len(ents))

	for _, ent := range ents {
		result = append(result, domain.AuthorizationRuleData{
			ID:        ent.ID.String(),
			OwnerID:   ent.OwnerID.String(),
			ViewerID:  ent.ViewerID.String(),
			Resource:  ent.Resource,
			Effect:    ent.Effect,
			CreatedAt: ent.CreatedAt,
		})
	}

	return result
}
