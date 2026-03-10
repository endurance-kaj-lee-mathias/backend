package application

import (
	"strconv"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) decryptStressSamples(ents []entities.StressSampleExportEntity, userKey []byte) ([]domain.StressSampleData, error) {
	result := make([]domain.StressSampleData, 0, len(ents))

	for _, ent := range ents {
		meanHR, err := s.decryptFloat(ent.EncryptedMeanHR, userKey)
		if err != nil {
			return nil, err
		}

		rmssd, err := s.decryptFloat(ent.EncryptedRMSSDms, userKey)
		if err != nil {
			return nil, err
		}

		var restingHR *float64
		if len(ent.EncryptedRestingHR) > 0 {
			val, err := s.decryptFloat(ent.EncryptedRestingHR, userKey)
			if err != nil {
				return nil, err
			}
			restingHR = &val
		}

		var steps *int
		if len(ent.EncryptedSteps) > 0 {
			val, err := s.decryptInt(ent.EncryptedSteps, userKey)
			if err != nil {
				return nil, err
			}
			steps = &val
		}

		var sleepDebt *float64
		if len(ent.EncryptedSleepDebtHours) > 0 {
			val, err := s.decryptFloat(ent.EncryptedSleepDebtHours, userKey)
			if err != nil {
				return nil, err
			}
			sleepDebt = &val
		}

		result = append(result, domain.StressSampleData{
			ID:             ent.ID.String(),
			TimestampUTC:   ent.TimestampUTC,
			WindowMinutes:  ent.WindowMinutes,
			MeanHR:         meanHR,
			RMSSDms:        rmssd,
			RestingHR:      restingHR,
			Steps:          steps,
			SleepDebtHours: sleepDebt,
			CreatedAt:      ent.CreatedAt,
		})
	}

	return result, nil
}

func (s *service) mapStressScores(ents []entities.StressScoreExportEntity) []domain.StressScoreData {
	result := make([]domain.StressScoreData, 0, len(ents))

	for _, ent := range ents {
		result = append(result, domain.StressScoreData{
			ID:           ent.ID.String(),
			Score:        ent.Score,
			Category:     ent.Category,
			ModelVersion: ent.ModelVersion,
			ComputedAt:   ent.ComputedAt,
		})
	}

	return result
}

func (s *service) decryptMoodEntries(ents []entities.MoodEntryExportEntity, userKey []byte) ([]domain.MoodEntryData, error) {
	result := make([]domain.MoodEntryData, 0, len(ents))

	for _, ent := range ents {
		var notes *string
		if len(ent.EncryptedNotes) > 0 {
			decrypted, err := s.enc.Decrypt(ent.EncryptedNotes, userKey)
			if err != nil {
				return nil, err
			}
			n := string(decrypted)
			notes = &n
		}

		result = append(result, domain.MoodEntryData{
			ID:        ent.ID.String(),
			Date:      ent.Date,
			MoodScore: ent.MoodScore,
			Notes:     notes,
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
		})
	}

	return result, nil
}

func (s *service) decryptFloat(ciphertext []byte, key []byte) (float64, error) {
	plaintext, err := s.enc.Decrypt(ciphertext, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(plaintext), 64)
}

func (s *service) decryptInt(ciphertext []byte, key []byte) (int, error) {
	plaintext, err := s.enc.Decrypt(ciphertext, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(plaintext))
}
