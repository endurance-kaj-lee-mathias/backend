package application

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) decryptSupportMembers(ents []entities.SupportMemberExportEntity) ([]domain.SupportMemberData, error) {
	result := make([]domain.SupportMemberData, 0, len(ents))

	for _, ent := range ents {
		memberKey, err := s.enc.DecryptUserKey(ent.EncryptedUserKey)
		if err != nil {
			return nil, err
		}

		email, err := s.enc.Decrypt(ent.EncryptedEmail, memberKey)
		if err != nil {
			return nil, err
		}

		username, err := s.enc.Decrypt(ent.EncryptedUsername, memberKey)
		if err != nil {
			return nil, err
		}

		firstName, err := s.enc.Decrypt(ent.EncryptedFirst, memberKey)
		if err != nil {
			return nil, err
		}

		lastName, err := s.enc.Decrypt(ent.EncryptedLast, memberKey)
		if err != nil {
			return nil, err
		}

		result = append(result, domain.SupportMemberData{
			ID:        ent.ID.String(),
			Email:     string(email),
			Username:  string(username),
			FirstName: string(firstName),
			LastName:  string(lastName),
			Image:     ent.Image,
			CreatedAt: ent.CreatedAt,
		})
	}

	return result, nil
}
