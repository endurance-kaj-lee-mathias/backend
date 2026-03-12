package models

import "time"

type StressSampleRequest struct {
	UserID         string    `json:"userId"`
	TimestampUTC   time.Time `json:"timestamp"`
	WindowMinutes  int       `json:"windowMinutes"`
	MeanHR         float64   `json:"meanHr"`
	RMSSDms        float64   `json:"rmssdMs"`
	RestingHR      *float64  `json:"restingHr,omitempty"`
	Steps          *int      `json:"steps,omitempty"`
	SleepDebtHours *float64  `json:"sleepDebtHours,omitempty"`
}

func (m *StressSampleRequest) Validate() error {
	if m.UserID == "" {
		return InvalidUserID
	}
	if m.TimestampUTC.IsZero() {
		return InvalidTimestamp
	}
	if m.WindowMinutes <= 0 {
		return InvalidWindowMinutes
	}
	if m.MeanHR <= 0 {
		return InvalidMeanHR
	}
	if m.RMSSDms <= 0 {
		return InvalidRMSSD
	}
	return nil
}

type LatestSampleResponse struct {
	Timestamp time.Time `json:"timestamp"`
}
