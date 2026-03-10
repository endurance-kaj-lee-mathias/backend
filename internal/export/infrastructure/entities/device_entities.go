package entities

import (
	"time"
)

type DeviceExportEntity struct {
	DeviceToken string
	Platform    string
	CreatedAt   time.Time
}
