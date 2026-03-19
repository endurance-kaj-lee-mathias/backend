package application

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	authzdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/journal/infrastructure/entities"
)

func (s *service) GetJournal(ctx context.Context, viewerID uuid.UUID, veteranID uuid.UUID, weekOffset int) (domain.JournalReport, error) {
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

	profileEnt, err := s.repo.GetUserProfile(ctx, veteranID)
	if err != nil {
		return domain.JournalReport{}, err
	}

	userKey, err := s.enc.DecryptUserKey(profileEnt.EncryptedUserKey)
	if err != nil {
		return domain.JournalReport{}, err
	}

	rolesBytes, err := s.enc.Decrypt(profileEnt.EncryptedRoles, userKey)
	if err != nil {
		return domain.JournalReport{}, err
	}

	var roles []string
	if err := json.Unmarshal(rolesBytes, &roles); err != nil {
		return domain.JournalReport{}, err
	}

	isVeteran := false
	for _, r := range roles {
		if r == "veteran" {
			isVeteran = true
			break
		}
	}

	if !isVeteran {
		return domain.JournalReport{}, domain.NotVeteran
	}

	if profileAllowed {
		profile, err := s.decryptProfile(profileEnt, userKey)
		if err != nil {
			return domain.JournalReport{}, err
		}
		report.UserProfile = &profile
	}

	if moodAllowed {
		rows, err := s.repo.GetWeeklyAverages(ctx, veteranID, weekOffset)
		if err != nil {
			return domain.JournalReport{}, err
		}

		noteRows, err := s.repo.GetWeeklyMoodNotes(ctx, veteranID, weekOffset)
		if err != nil {
			return domain.JournalReport{}, err
		}

		notesByDate, err := decryptNotesByDate(noteRows, s.enc, userKey)
		if err != nil {
			return domain.JournalReport{}, err
		}

		total := totalFromRows(rows)
		days := toDailyAverages(rows, stressAllowed, notesByDate)
		report.Weekly = &domain.WeeklyPage{Days: days, Total: total}
	}

	return report, nil
}

func totalFromRows(rows []entities.DailyAverageRow) int {
	if len(rows) == 0 {
		return 0
	}
	return rows[0].Total
}

func toDailyAverages(rows []entities.DailyAverageRow, stressAllowed bool, notesByDate map[string][]string) []domain.DailyAverage {
	days := make([]domain.DailyAverage, 0, len(rows))

	for _, row := range rows {
		avgStress := row.AvgStress
		if !stressAllowed {
			avgStress = nil
		}

		dateKey := row.Date.Format("2006-01-02")
		notes := notesByDate[dateKey]

		days = append(days, domain.DailyAverage{
			Date:      dateKey,
			AvgMood:   row.AvgMood,
			AvgStress: avgStress,
			Notes:     notes,
		})
	}

	return days
}

func decryptNotesByDate(rows []entities.MoodEntryNoteRow, enc encryption.Service, userKey []byte) (map[string][]string, error) {
	result := make(map[string][]string)

	for _, row := range rows {
		if len(row.EncryptedNotes) == 0 {
			continue
		}

		decrypted, err := enc.Decrypt(row.EncryptedNotes, userKey)
		if err != nil {
			return nil, err
		}

		dateKey := row.Date.Format("2006-01-02")
		result[dateKey] = append(result[dateKey], string(decrypted))
	}

	return result, nil
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
