package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type VeteranMoodSummary struct {
	VeteranID     uuid.UUID
	FirstName     string
	LastName      string
	Image         string
	LatestScore   *int
	LastUpdatedAt *time.Time
}
