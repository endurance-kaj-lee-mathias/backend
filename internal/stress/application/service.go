package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure/entities"
)

func (s *service) IngestSample(ctx context.Context, sample domain.StressSample) error {
	encryptedUserKey, err := s.userKeyReader.GetEncryptedUserKey(ctx, sample.UserID)
	if err != nil {
		return err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return err
	}

	ent, err := entities.ToEntity(sample, s.enc, userKey)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, ent)
}
