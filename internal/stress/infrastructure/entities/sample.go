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

func FromEntity(ent StressSampleEntity, enc encryption.Service, userKey []byte) (domain.StressSample, error) {
	meanHRBytes, err := enc.Decrypt(ent.EncryptedMeanHR, userKey)
	if err != nil {
		return domain.StressSample{}, err
	}
	meanHR, err := strconv.ParseFloat(string(meanHRBytes), 64)
	if err != nil {
		return domain.StressSample{}, err
	}

	rmssdBytes, err := enc.Decrypt(ent.EncryptedRMSSDms, userKey)
	if err != nil {
		return domain.StressSample{}, err
	}
	rmssd, err := strconv.ParseFloat(string(rmssdBytes), 64)
	if err != nil {
		return domain.StressSample{}, err
	}

	var restingHR *float64
	if len(ent.EncryptedRestingHR) > 0 {
		b, err := enc.Decrypt(ent.EncryptedRestingHR, userKey)
		if err != nil {
			return domain.StressSample{}, err
		}
		v, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return domain.StressSample{}, err
		}
		restingHR = &v
	}

	var steps *int
	if len(ent.EncryptedSteps) > 0 {
		b, err := enc.Decrypt(ent.EncryptedSteps, userKey)
		if err != nil {
			return domain.StressSample{}, err
		}
		v, err := strconv.Atoi(string(b))
		if err != nil {
			return domain.StressSample{}, err
		}
		steps = &v
	}

	var sleepDebt *float64
	if len(ent.EncryptedSleepDebtHours) > 0 {
		b, err := enc.Decrypt(ent.EncryptedSleepDebtHours, userKey)
		if err != nil {
			return domain.StressSample{}, err
		}
		v, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return domain.StressSample{}, err
		}
		sleepDebt = &v
	}

	return domain.StressSample{
		ID:             domain.SampleId{UUID: ent.ID},
		UserID:         ent.UserID,
		TimestampUTC:   ent.TimestampUTC,
		WindowMinutes:  ent.WindowMinutes,
		MeanHR:         meanHR,
		RMSSDms:        rmssd,
		RestingHR:      restingHR,
		Steps:          steps,
		SleepDebtHours: sleepDebt,
		CreatedAt:      ent.CreatedAt,
	}, nil
}
