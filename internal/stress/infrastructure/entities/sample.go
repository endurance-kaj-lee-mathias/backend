package entities

import (
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
)

type StressSampleEntity struct {
	ID                      uuid.UUID `db:"id"`
	UserID                  uuid.UUID `db:"user_id"`
	TimestampUTC            time.Time `db:"timestamp_utc"`
	WindowMinutes           int       `db:"window_minutes"`
	EncryptedMeanHR         []byte    `db:"encrypted_mean_hr"`
	EncryptedRMSSDms        []byte    `db:"encrypted_rmssd_ms"`
	EncryptedRestingHR      []byte    `db:"encrypted_resting_hr"`
	EncryptedSteps          []byte    `db:"encrypted_steps"`
	EncryptedSleepDebtHours []byte    `db:"encrypted_sleep_debt_hours"`
	CreatedAt               time.Time `db:"created_at"`
}

func ToEntity(sample domain.StressSample, enc encryption.Service, userKey []byte) (StressSampleEntity, error) {
	encMeanHR, err := enc.Encrypt([]byte(strconv.FormatFloat(sample.MeanHR, 'f', -1, 64)), userKey)
	if err != nil {
		return StressSampleEntity{}, err
	}

	encRMSSD, err := enc.Encrypt([]byte(strconv.FormatFloat(sample.RMSSDms, 'f', -1, 64)), userKey)
	if err != nil {
		return StressSampleEntity{}, err
	}

	var encRestingHR []byte
	if sample.RestingHR != nil {
		encRestingHR, err = enc.Encrypt([]byte(strconv.FormatFloat(*sample.RestingHR, 'f', -1, 64)), userKey)
		if err != nil {
			return StressSampleEntity{}, err
		}
	}

	var encSteps []byte
	if sample.Steps != nil {
		encSteps, err = enc.Encrypt([]byte(strconv.Itoa(*sample.Steps)), userKey)
		if err != nil {
			return StressSampleEntity{}, err
		}
	}

	var encSleepDebt []byte
	if sample.SleepDebtHours != nil {
		encSleepDebt, err = enc.Encrypt([]byte(strconv.FormatFloat(*sample.SleepDebtHours, 'f', -1, 64)), userKey)
		if err != nil {
			return StressSampleEntity{}, err
		}
	}

	return StressSampleEntity{
		ID:                      sample.ID.UUID,
		UserID:                  sample.UserID,
		TimestampUTC:            sample.TimestampUTC,
		WindowMinutes:           sample.WindowMinutes,
		EncryptedMeanHR:         encMeanHR,
		EncryptedRMSSDms:        encRMSSD,
		EncryptedRestingHR:      encRestingHR,
		EncryptedSteps:          encSteps,
		EncryptedSleepDebtHours: encSleepDebt,
		CreatedAt:               sample.CreatedAt,
	}, nil
}
