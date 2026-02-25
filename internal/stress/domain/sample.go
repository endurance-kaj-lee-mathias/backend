package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type StressSample struct {
	ID             SampleId
	UserID         uuid.UUID
	TimestampUTC   time.Time
	WindowMinutes  int
	MeanHR         float64
	RMSSDms        float64
	RestingHR      *float64
	Steps          *int
	SleepDebtHours *float64
	CreatedAt      time.Time
}

func NewStressSample(
	userID uuid.UUID,
	timestampUTC time.Time,
	windowMinutes int,
	meanHR float64,
	rmssdMs float64,
	restingHR *float64,
	steps *int,
	sleepDebtHours *float64,
) (StressSample, error) {
	id, err := NewSampleId()
	if err != nil {
		return StressSample{}, err
	}

	return StressSample{
		ID:             id,
		UserID:         userID,
		TimestampUTC:   timestampUTC,
		WindowMinutes:  windowMinutes,
		MeanHR:         meanHR,
		RMSSDms:        rmssdMs,
		RestingHR:      restingHR,
		Steps:          steps,
		SleepDebtHours: sleepDebtHours,
		CreatedAt:      time.Now().UTC(),
	}, nil
}
