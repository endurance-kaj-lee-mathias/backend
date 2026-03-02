package models

import "time"

type CreateSlotRequest struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	IsUrgent  bool      `json:"isUrgent"`
}

func (m *CreateSlotRequest) Validate() error {
	if m.StartTime.IsZero() {
		return InvalidStartTime
	}
	if m.EndTime.IsZero() {
		return InvalidEndTime
	}
	return nil
}
