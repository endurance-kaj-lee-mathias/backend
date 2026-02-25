package models

import "time"

type StressSampleRequest struct {
	UserID         string    `json:"user_id"`
	TimestampUTC   time.Time `json:"timestamp_utc"`
	WindowMinutes  int       `json:"window_minutes"`
	MeanHR         float64   `json:"mean_hr"`
	RMSSDms        float64   `json:"rmssd_ms"`
	RestingHR      *float64  `json:"resting_hr,omitempty"`
	Steps          *int      `json:"steps,omitempty"`
	SleepDebtHours *float64  `json:"sleep_debt_hours,omitempty"`
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
