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

func (s *service) GetEntriesByUserID(ctx context.Context, userID domain.UserId) ([]domain.MoodEntry, error) {
	encryptedUserKey, err := s.userKeyReader.GetEncryptedUserKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return nil, err
	}

	rawEntries, err := s.repo.FindAllByUserID(ctx, userID.UUID)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.MoodEntry, 0, len(rawEntries))

	for _, ent := range rawEntries {
		entry, err := entities.FromEntity(ent, s.enc, userKey)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
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

func (s *service) GetVeteransMood(ctx context.Context, memberID uuid.UUID) ([]VeteranMoodSummary, error) {
	veterans, err := s.veteranLister.GetVeteransForMember(ctx, memberID)
	if err != nil {
		return nil, err
	}

	summaries := make([]VeteranMoodSummary, 0, len(veterans))

	for _, veteran := range veterans {
		allowed, err := s.authz.IsAllowed(ctx, veteran.ID, memberID, "moodEntries")
		if err != nil {
			return nil, err
		}

		if !allowed {
			continue
		}

		entries, err := s.GetEntriesByUserID(ctx, domain.UserId{UUID: veteran.ID})
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, VeteranMoodSummary{
			VeteranID: veteran.ID,
			FirstName: veteran.FirstName,
			LastName:  veteran.LastName,
			Image:     veteran.Image,
			Entries:   entries,
		})
	}

	return summaries, nil
}
