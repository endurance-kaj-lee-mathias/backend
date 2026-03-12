package application

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure/entities"
)

const minSamplesForComputation = 12

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

	if err := s.repo.Create(ctx, ent); err != nil {
		return err
	}

	count, err := s.repo.CountSamples(ctx, sample.UserID)
	if err != nil {
		slog.Error("stress: count samples after ingest", "userID", sample.UserID, "error", err)
		return nil
	}

	if count > minSamplesForComputation {
		if _, err := s.computeStressScore(ctx, sample.UserID); err != nil {
			slog.Error("stress: auto-compute score after ingest", "userID", sample.UserID, "error", err)
		}
	}

	return nil
}

func (s *service) computeStressScore(ctx context.Context, userID uuid.UUID) (domain.StressScore, error) {
	encryptedUserKey, err := s.userKeyReader.GetEncryptedUserKey(ctx, userID)
	if err != nil {
		return domain.StressScore{}, err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return domain.StressScore{}, err
	}

	encryptedSamples, err := s.repo.GetSamplesLast90Days(ctx, userID)
	if err != nil {
		return domain.StressScore{}, err
	}

	samples := make([]domain.StressSample, 0, len(encryptedSamples))
	for _, ent := range encryptedSamples {
		sample, err := entities.FromEntity(ent, s.enc, userKey)
		if err != nil {
			slog.Error("stress: decrypt sample", "userID", userID, "sampleID", ent.ID, "error", err)
			return domain.StressScore{}, err
		}
		samples = append(samples, sample)
	}

	score, err := s.algoClient.ComputeScore(ctx, userID, samples)
	if err != nil {
		return domain.StressScore{}, err
	}

	scoreEnt := entities.ScoreToEntity(score)
	if err := s.repo.CreateScore(ctx, scoreEnt); err != nil {
		slog.Error("stress: persist score", "userID", userID, "error", err)
		return domain.StressScore{}, err
	}

	return score, nil
}

func (s *service) GetLatestScore(ctx context.Context, userID uuid.UUID) (domain.StressScore, error) {
	return s.repo.GetLatestScore(ctx, userID)
}

func (s *service) GetLatestSampleTimestamp(ctx context.Context, userID uuid.UUID) (time.Time, error) {
	return s.repo.GetLatestSampleTimestamp(ctx, userID)
}

func (s *service) DeleteMySamples(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteAllByUserID(ctx, userID)
}

func (s *service) GetScoresPaginated(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.StressScore, int, error) {
	return s.repo.GetScoresPaginated(ctx, userID, limit, offset)
}
