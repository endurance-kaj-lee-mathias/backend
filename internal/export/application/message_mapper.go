package application

import (
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) decryptMessages(ents []entities.MessageExportEntity, userKey []byte) ([]domain.MessageData, error) {
	result := make([]domain.MessageData, 0, len(ents))
	convKeyCache := make(map[uuid.UUID][]byte)
	otherUserCache := make(map[uuid.UUID][3]string)

	for _, ent := range ents {
		convKey, ok := convKeyCache[ent.ConversationID]
		if !ok {
			decryptedConvKey, err := s.enc.Decrypt(ent.EncryptedConversationKey, userKey)
			if err != nil {
				return nil, err
			}
			convKey = decryptedConvKey
			convKeyCache[ent.ConversationID] = convKey
		}

		content, err := s.enc.Decrypt(ent.EncryptedContent, convKey)
		if err != nil {
			return nil, err
		}

		otherInfo, ok := otherUserCache[ent.OtherParticipantID]
		if !ok {
			otherUserKey, err := s.enc.DecryptUserKey(ent.OtherParticipantEncryptedUserKey)
			if err != nil {
				return nil, err
			}

			username, err := s.enc.Decrypt(ent.OtherParticipantEncryptedUsername, otherUserKey)
			if err != nil {
				return nil, err
			}

			firstName, err := s.enc.Decrypt(ent.OtherParticipantEncryptedFirstName, otherUserKey)
			if err != nil {
				return nil, err
			}

			lastName, err := s.enc.Decrypt(ent.OtherParticipantEncryptedLastName, otherUserKey)
			if err != nil {
				return nil, err
			}

			otherInfo = [3]string{string(username), string(firstName), string(lastName)}
			otherUserCache[ent.OtherParticipantID] = otherInfo
		}

		result = append(result, domain.MessageData{
			ID:                        ent.ID.String(),
			ConversationID:            ent.ConversationID.String(),
			SenderID:                  ent.SenderID.String(),
			Content:                   string(content),
			CreatedAt:                 ent.CreatedAt,
			OtherParticipantUsername:  otherInfo[0],
			OtherParticipantFirstName: otherInfo[1],
			OtherParticipantLastName:  otherInfo[2],
		})
	}

	return result, nil
}
