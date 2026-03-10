package application

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) decryptInvites(ents []entities.InviteExportEntity) ([]domain.InviteData, error) {
	result := make([]domain.InviteData, 0, len(ents))

	for _, ent := range ents {
		otherKey, err := s.enc.DecryptUserKey(ent.OtherEncryptedUserKey)
		if err != nil {
			return nil, err
		}

		username, err := s.enc.Decrypt(ent.OtherEncryptedUsername, otherKey)
		if err != nil {
			return nil, err
		}

		firstName, err := s.enc.Decrypt(ent.OtherEncryptedFirstName, otherKey)
		if err != nil {
			return nil, err
		}

		lastName, err := s.enc.Decrypt(ent.OtherEncryptedLastName, otherKey)
		if err != nil {
			return nil, err
		}

		otherImage := ""
		if ent.OtherImage != nil {
			otherImage = *ent.OtherImage
		}

		result = append(result, domain.InviteData{
			ID:             ent.ID.String(),
			OtherUserID:    ent.OtherUserID.String(),
			OtherUsername:  string(username),
			OtherFirstName: string(firstName),
			OtherLastName:  string(lastName),
			OtherImage:     otherImage,
			Status:         ent.Status,
			Note:           ent.Note,
			CreatedAt:      ent.CreatedAt,
			UpdatedAt:      ent.UpdatedAt,
		})
	}

	return result, nil
}
