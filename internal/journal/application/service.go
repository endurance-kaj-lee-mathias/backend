package application

import (
	"context"

	"github.com/gofrs/uuid"
	authzdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/infrastructure/entities"
)

func (s *service) GetJournal(ctx context.Context, viewerID uuid.UUID, veteranID uuid.UUID, limit, offset int) (domain.JournalReport, error) {
	report := domain.JournalReport{VeteranID: veteranID}

	profileAllowed, err := s.authz.IsAllowed(ctx, veteranID, viewerID, string(authzdomain.ResourceUserProfile))
	if err != nil {
		return domain.JournalReport{}, err
	}

	stressAllowed, err := s.authz.IsAllowed(ctx, veteranID, viewerID, string(authzdomain.ResourceStressScores))
	if err != nil {
		return domain.JournalReport{}, err
	}

	moodAllowed, err := s.authz.IsAllowed(ctx, veteranID, viewerID, string(authzdomain.ResourceMoodEntries))
	if err != nil {
		return domain.JournalReport{}, err
	}

	var userKey []byte

	if profileAllowed || moodAllowed {
		profileEnt, err := s.repo.GetUserProfile(ctx, veteranID)
		if err != nil {
			return domain.JournalReport{}, err
		}

		userKey, err = s.enc.DecryptUserKey(profileEnt.EncryptedUserKey)
		if err != nil {
			return domain.JournalReport{}, err
		}

		if profileAllowed {
			profile, err := s.decryptProfile(profileEnt, userKey)
			if err != nil {
				return domain.JournalReport{}, err
			}
			report.UserProfile = &profile
		}
	}

	if stressAllowed {
		rows, total, err := s.repo.GetStressScoresPaginated(ctx, veteranID, limit, offset)
		if err != nil {
			return domain.JournalReport{}, err
		}

		items := make([]domain.StressScoreItem, 0, len(rows))
		for _, row := range rows {
			items = append(items, domain.StressScoreItem{
				ID:           row.ID,
				Score:        row.Score,
				Category:     row.Category,
				ModelVersion: row.ModelVersion,
				ComputedAt:   row.ComputedAt,
			})
		}

		report.StressScores = &domain.ScoresPage{Items: items, Total: total}
	}

	if moodAllowed {
		rows, total, err := s.repo.GetMoodEntriesPaginated(ctx, veteranID, limit, offset)
		if err != nil {
			return domain.JournalReport{}, err
		}

		items := make([]domain.MoodEntryItem, 0, len(rows))
		for _, row := range rows {
			var notes *string

			if len(row.EncryptedNotes) > 0 {
				decrypted, err := s.enc.Decrypt(row.EncryptedNotes, userKey)
				if err != nil {
					return domain.JournalReport{}, err
				}
				n := string(decrypted)
				notes = &n
			}

			items = append(items, domain.MoodEntryItem{
				ID:        row.ID,
				Date:      row.Date,
				MoodScore: row.MoodScore,
				Notes:     notes,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			})
		}

		report.MoodEntries = &domain.MoodPage{Items: items, Total: total}
	}

	return report, nil
}

func (s *service) decryptProfile(ent entities.UserProfileEntity, userKey []byte) (domain.UserProfileSection, error) {
	firstName, err := s.enc.Decrypt(ent.EncryptedFirstName, userKey)
	if err != nil {
		return domain.UserProfileSection{}, err
	}

	lastName, err := s.enc.Decrypt(ent.EncryptedLastName, userKey)
	if err != nil {
		return domain.UserProfileSection{}, err
	}

	username, err := s.enc.Decrypt(ent.EncryptedUsername, userKey)
	if err != nil {
		return domain.UserProfileSection{}, err
	}

	about, err := s.enc.Decrypt(ent.EncryptedAbout, userKey)
	if err != nil {
		return domain.UserProfileSection{}, err
	}

	introduction, err := s.enc.Decrypt(ent.EncryptedIntroduction, userKey)
	if err != nil {
		return domain.UserProfileSection{}, err
	}

	var phoneNumber *string

	if len(ent.EncryptedPhoneNumber) > 0 {
		decrypted, err := s.enc.Decrypt(ent.EncryptedPhoneNumber, userKey)
		if err != nil {
			return domain.UserProfileSection{}, err
		}
		n := string(decrypted)
		phoneNumber = &n
	}

	return domain.UserProfileSection{
		FirstName:    string(firstName),
		LastName:     string(lastName),
		Username:     string(username),
		About:        string(about),
		Introduction: string(introduction),
		Image:        ent.Image,
		PhoneNumber:  phoneNumber,
		IsPrivate:    ent.IsPrivate,
	}, nil
}
