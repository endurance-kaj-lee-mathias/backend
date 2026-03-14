package application

import (
	"context"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure"
)

type Service interface {
	UpsertMoodEntry(ctx context.Context, entry domain.MoodEntry) error
	GetEntryByID(ctx context.Context, id domain.MoodId) (*domain.MoodEntry, error)
	GetEntriesByUserID(ctx context.Context, userID domain.UserId) ([]domain.MoodEntry, error)
	GetTodayEntry(ctx context.Context, userID domain.UserId) (*domain.MoodEntry, error)
	UpdateMoodEntry(ctx context.Context, entry domain.MoodEntry) error
	DeleteMoodEntry(ctx context.Context, id domain.MoodId) error
	DeleteMyMoodEntries(ctx context.Context, userID domain.UserId) error
	GetVeteransSupport(ctx context.Context, memberID uuid.UUID) ([]domain.VeteranMoodSummary, error)
}

type AuthorizationChecker interface {
	IsAllowed(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID, resource string) (bool, error)
}

type PhoneNotifier interface {
	Notify(ctx context.Context, deviceToken string) error
}

type service struct {
	repo          infrastructure.Repository
	userKeyReader infrastructure.UserKeyReader
	enc           encryption.Service
	veteranLister infrastructure.VeteranLister
	authz         AuthorizationChecker
}

func NewService(repo infrastructure.Repository, userKeyReader infrastructure.UserKeyReader, enc encryption.Service, veteranLister infrastructure.VeteranLister, authz AuthorizationChecker) Service {
	return &service{repo: repo, userKeyReader: userKeyReader, enc: enc, veteranLister: veteranLister, authz: authz}
}
