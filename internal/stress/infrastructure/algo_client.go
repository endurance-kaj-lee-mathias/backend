package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
)

type algoSample struct {
	TimestampUTC   time.Time `json:"timestamp_utc"`
	WindowMinutes  int       `json:"window_minutes"`
	MeanHR         float64   `json:"mean_hr"`
	RMSSDms        float64   `json:"rmssd_ms"`
	RestingHR      *float64  `json:"resting_hr,omitempty"`
	Steps          *int      `json:"steps,omitempty"`
	SleepDebtHours *float64  `json:"sleep_debt_hours,omitempty"`
}

type algoRequest struct {
	UserID  uuid.UUID    `json:"user_id"`
	Samples []algoSample `json:"samples"`
}

type algoResponse struct {
	Score        float64  `json:"score"`
	Category     string   `json:"category"`
	ModelVersion string   `json:"model_version"`
	ZHR          *float64 `json:"z_hr,omitempty"`
	ZRMSSD       *float64 `json:"z_rmssd,omitempty"`
}

func (c *algoClient) ComputeScore(ctx context.Context, userID uuid.UUID, samples []domain.StressSample) (domain.StressScore, error) {
	payload := algoRequest{
		UserID:  userID,
		Samples: make([]algoSample, len(samples)),
	}

	for i, s := range samples {
		payload.Samples[i] = algoSample{
			TimestampUTC:   s.TimestampUTC,
			WindowMinutes:  s.WindowMinutes,
			MeanHR:         s.MeanHR,
			RMSSDms:        s.RMSSDms,
			RestingHR:      s.RestingHR,
			Steps:          s.Steps,
			SleepDebtHours: s.SleepDebtHours,
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return domain.StressScore{}, fmt.Errorf("algo: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/stress/compute", bytes.NewReader(body))
	if err != nil {
		return domain.StressScore{}, fmt.Errorf("algo: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("algo: request failed", "userID", userID, "error", err)
		return domain.StressScore{}, AlgoServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		slog.Error("algo: unexpected status", "userID", userID, "status", resp.StatusCode, "body", string(raw))
		return domain.StressScore{}, AlgoServiceError
	}

	var result algoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("algo: decode response", "userID", userID, "error", err)
		return domain.StressScore{}, fmt.Errorf("algo: decode response: %w", err)
	}

	return domain.NewStressScore(userID, result.Score, result.Category, result.ModelVersion)
}
