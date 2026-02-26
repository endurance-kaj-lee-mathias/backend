package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure/entities"
)

func (s *service) UpsertMoodEntry(ctx context.Context, entry domain.MoodEntry) error {
	encryptedUserKey, err := s.userKeyReader.GetEncryptedUserKey(ctx, entry.UserID)
	if err != nil {
		return err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return err
	}

	ent, err := entities.ToEntity(entry, s.enc, userKey)
	if err != nil {
		return err
	}

	return s.repo.Upsert(ctx, ent)
}
