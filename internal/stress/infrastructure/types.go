package infrastructure

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/infrastructure/entities"
)

type Repository interface {
	Create(ctx context.Context, ent entities.StressSampleEntity) error
	CountSamples(ctx context.Context, userID uuid.UUID) (int, error)
	GetLatestSampleTimestamp(ctx context.Context, userID uuid.UUID) (time.Time, error)
	GetSamplesLast90Days(ctx context.Context, userID uuid.UUID) ([]entities.StressSampleEntity, error)
	CreateScore(ctx context.Context, ent entities.StressScoreEntity) error
	GetLatestScore(ctx context.Context, userID uuid.UUID) (domain.StressScore, error)
	GetScoresPaginated(ctx context.Context, userID uuid.UUID, weekOffset int) ([]domain.StressScore, int, error)
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
}

type UserKeyReader interface {
	GetEncryptedUserKey(ctx context.Context, userID uuid.UUID) ([]byte, error)
}

type AlgoClient interface {
	ComputeScore(ctx context.Context, userID uuid.UUID, samples []domain.StressSample) (domain.StressScore, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

type userKeyReader struct {
	db  *sql.DB
	enc encryption.Service
}

func NewUserKeyReader(db *sql.DB, enc encryption.Service) UserKeyReader {
	return &userKeyReader{db: db, enc: enc}
}

type algoClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewAlgoClient(baseURL string, apiKey string, httpClient *http.Client) AlgoClient {
	return &algoClient{baseURL: baseURL, apiKey: apiKey, httpClient: httpClient}
}
