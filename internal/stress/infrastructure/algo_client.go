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
	HeartRate        float64   `json:"heart_rate"`
	RMSSD            float64   `json:"rmssd"`
	RestingHeartRate float64   `json:"resting_heart_rate,omitempty"`
	Steps            float64   `json:"steps,omitempty"`
	SleepDebtHours   float64   `json:"sleep_debt_hours,omitempty"`
	RecordedAt       time.Time `json:"recorded_at"`
}

type algoRequest struct {
	UserID  uuid.UUID    `json:"user_id"`
	Samples []algoSample `json:"samples"`
}

type algoResponse struct {
	Score    float64 `json:"score"`
	Category string  `json:"category"`
	ZHR      float64 `json:"z_hr"`
	ZRMSSD   float64 `json:"z_rmssd"`
}

func (c *algoClient) ComputeScore(ctx context.Context, userID uuid.UUID, samples []domain.StressSample) (domain.StressScore, error) {
	const maxSamples = 30
	if len(samples) > maxSamples {
		samples = samples[len(samples)-maxSamples:]
	}

	payload := algoRequest{
		UserID:  userID,
		Samples: make([]algoSample, len(samples)),
	}

	for i, s := range samples {
		var restingHR float64
		if s.RestingHR != nil {
			restingHR = *s.RestingHR
		}
		var steps float64
		if s.Steps != nil {
			steps = float64(*s.Steps)
		}
		var sleepDebt float64
		if s.SleepDebtHours != nil {
			sleepDebt = *s.SleepDebtHours
		}
		payload.Samples[i] = algoSample{
			HeartRate:        s.MeanHR,
			RMSSD:            s.RMSSDms,
			RestingHeartRate: restingHR,
			Steps:            steps,
			SleepDebtHours:   sleepDebt,
			RecordedAt:       s.TimestampUTC,
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
	req.Header.Set("X-API-Key", c.apiKey)

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

	return domain.NewStressScore(userID, result.Score, result.Category, "")
}
