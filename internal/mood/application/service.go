package application

import (
	"context"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure/entities"
)

func (s *service) GetEntryByID(ctx context.Context, id domain.MoodId) (*domain.MoodEntry, error) {
	ent, err := s.repo.FindByID(ctx, id.UUID)
	if err != nil {
		return nil, err
	}

	encryptedUserKey, err := s.userKeyReader.GetEncryptedUserKey(ctx, domain.UserId{UUID: ent.UserID})
	if err != nil {
		return nil, err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return nil, err
	}

	entry, err := entities.FromEntity(*ent, s.enc, userKey)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (s *service) UpdateMoodEntry(ctx context.Context, entry domain.MoodEntry) error {
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

	return s.repo.Update(ctx, ent)
}

func (s *service) DeleteMoodEntry(ctx context.Context, id domain.MoodId) error {
	return s.repo.Delete(ctx, id.UUID)
}

func (s *service) DeleteMyMoodEntries(ctx context.Context, userID domain.UserId) error {
	return s.repo.DeleteAllByUserID(ctx, userID.UUID)
}

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

func (s *service) GetEntriesByUserID(ctx context.Context, userID domain.UserId, weekOffset int) ([]domain.MoodEntry, int, error) {
	encryptedUserKey, err := s.userKeyReader.GetEncryptedUserKey(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return nil, 0, err
	}

	rawEntries, total, err := s.repo.FindPaginatedByUserID(ctx, userID.UUID, weekOffset)
	if err != nil {
		return nil, 0, err
	}

	entries := make([]domain.MoodEntry, 0, len(rawEntries))

	for _, ent := range rawEntries {
		entry, err := entities.FromEntity(ent, s.enc, userKey)
		if err != nil {
			return nil, 0, err
		}
		entries = append(entries, entry)
	}

	return entries, total, nil
}

func (s *service) GetTodayEntry(ctx context.Context, userID domain.UserId) (*domain.MoodEntry, error) {
	encryptedUserKey, err := s.userKeyReader.GetEncryptedUserKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return nil, err
	}

	ent, err := s.repo.FindTodayByUserID(ctx, userID.UUID)
	if err != nil {
		return nil, err
	}

	entry, err := entities.FromEntity(*ent, s.enc, userKey)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (s *service) GetVeteransSupport(ctx context.Context, memberID uuid.UUID) ([]domain.VeteranMoodSummary, error) {
	veterans, err := s.veteranLister.GetVeteransForMember(ctx, memberID)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.VeteranMoodSummary, 0, len(veterans))

	for _, veteran := range veterans {
		allowed, err := s.authz.IsAllowed(ctx, veteran.ID, memberID, "moodEntries")
		if err != nil {
			return nil, err
		}

		if !allowed {
			continue
		}

		latest, err := s.repo.FindLatestByUserID(ctx, veteran.ID)
		if err != nil {
			return nil, err
		}

		summary := domain.VeteranMoodSummary{
			VeteranID: veteran.ID,
			Username:  veteran.Username,
			FirstName: veteran.FirstName,
			LastName:  veteran.LastName,
			Image:     veteran.Image,
		}

		if latest != nil {
			summary.LatestScore = &latest.MoodScore
			summary.LastUpdatedAt = &latest.UpdatedAt
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
