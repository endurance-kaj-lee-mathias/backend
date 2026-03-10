package domain

import (
	"time"
)

type AccountSettingsData struct {
	IsPrivate bool         `json:"isPrivate"`
	Devices   []DeviceData `json:"devices"`
}

type DeviceData struct {
	Token     string    `json:"token"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"createdAt"`
}
