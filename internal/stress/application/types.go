package application

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure"
	usersdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserRiskService interface {
	UpdateRiskLevel(ctx context.Context, id usersdomain.UserId, riskLevel usersdomain.RiskLevel) error
	FindDeviceTokensByUserID(ctx context.Context, userID usersdomain.UserId) ([]string, error)
}

type Notifier interface {
	NotifyHighStress(ctx context.Context, userID uuid.UUID, deviceTokens []string) error
}

type Service interface {
	IngestSample(ctx context.Context, sample domain.StressSample) error
	GetLatestScore(ctx context.Context, userID uuid.UUID) (domain.StressScore, error)
	GetLatestSampleTimestamp(ctx context.Context, userID uuid.UUID) (time.Time, error)
	GetScoresPaginated(ctx context.Context, userID uuid.UUID, weekOffset int) ([]domain.StressScore, int, error)
	DeleteMySamples(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	repo          infrastructure.Repository
	userKeyReader infrastructure.UserKeyReader
	algoClient    infrastructure.AlgoClient
	enc           encryption.Service
	riskSvc       UserRiskService
	notifier      Notifier
}

func NewService(repo infrastructure.Repository, userKeyReader infrastructure.UserKeyReader, algoClient infrastructure.AlgoClient, enc encryption.Service, riskSvc UserRiskService, notifier Notifier) Service {
	return &service{repo: repo, userKeyReader: userKeyReader, algoClient: algoClient, enc: enc, riskSvc: riskSvc, notifier: notifier}
}
