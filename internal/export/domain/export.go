package domain

import (
	"time"
)

type UserDataExport struct {
	ExportedAt time.Time  `json:"exportedAt"`
	UserID     string     `json:"userId"`
	Data       ExportData `json:"data"`
}

type ExportData struct {
	Profile         ProfileData         `json:"profile"`
	Address         *AddressData        `json:"address,omitempty"`
	HealthData      HealthData          `json:"healthData"`
	Messages        []MessageData       `json:"messages"`
	Calendar        CalendarData        `json:"calendar"`
	SupportNetwork  SupportNetworkData  `json:"supportNetwork"`
	DataSharing     DataSharingData     `json:"dataSharing"`
	AccountSettings AccountSettingsData `json:"accountSettings"`
	Invites         InvitesData         `json:"invites"`
}
