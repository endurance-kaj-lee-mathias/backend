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

type ProfileData struct {
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	PhoneNumber  *string   `json:"phoneNumber,omitempty"`
	Roles        []string  `json:"roles"`
	About        string    `json:"about"`
	Introduction string    `json:"introduction"`
	Image        string    `json:"image"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AddressData struct {
	Street     string    `json:"street"`
	Locality   string    `json:"locality"`
	Region     string    `json:"region"`
	PostalCode string    `json:"postalCode"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"createdAt"`
}

type HealthData struct {
	StressSamples []StressSampleData `json:"stressSamples"`
	StressScores  []StressScoreData  `json:"stressScores"`
	MoodEntries   []MoodEntryData    `json:"moodEntries"`
}

type StressSampleData struct {
	ID             string    `json:"id"`
	TimestampUTC   time.Time `json:"timestampUtc"`
	WindowMinutes  int       `json:"windowMinutes"`
	MeanHR         float64   `json:"meanHr"`
	RMSSDms        float64   `json:"rmssdMs"`
	RestingHR      *float64  `json:"restingHr,omitempty"`
	Steps          *int      `json:"steps,omitempty"`
	SleepDebtHours *float64  `json:"sleepDebtHours,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type StressScoreData struct {
	ID           string    `json:"id"`
	Score        float64   `json:"score"`
	Category     string    `json:"category"`
	ModelVersion string    `json:"modelVersion"`
	ComputedAt   time.Time `json:"computedAt"`
}

type MoodEntryData struct {
	ID        string    `json:"id"`
	Date      time.Time `json:"date"`
	MoodScore int       `json:"moodScore"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MessageData struct {
	ID                        string    `json:"id"`
	ConversationID            string    `json:"conversationId"`
	SenderID                  string    `json:"senderId"`
	Content                   string    `json:"content"`
	CreatedAt                 time.Time `json:"createdAt"`
	OtherParticipantUsername  string    `json:"otherParticipantUsername"`
	OtherParticipantFirstName string    `json:"otherParticipantFirstName"`
	OtherParticipantLastName  string    `json:"otherParticipantLastName"`
}

type CalendarData struct {
	Appointments []AppointmentData `json:"appointments"`
	Slots        []SlotData        `json:"slots"`
}

type AppointmentData struct {
	ID               string    `json:"id"`
	SlotID           string    `json:"slotId"`
	VeteranID        string    `json:"veteranId"`
	ProviderID       string    `json:"providerId"`
	ProviderUsername string    `json:"providerUsername"`
	Status           string    `json:"status"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	IsUrgent         bool      `json:"isUrgent"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type SlotData struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	IsUrgent  bool      `json:"isUrgent"`
	IsBooked  bool      `json:"isBooked"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SupportNetworkData struct {
	Supporters        []SupportMemberData `json:"supporters"`
	SupportedVeterans []SupportMemberData `json:"supportedVeterans"`
}

type SupportMemberData struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

type DataSharingData struct {
	SharedByMe   []AuthorizationRuleData `json:"sharedByMe"`
	SharedWithMe []AuthorizationRuleData `json:"sharedWithMe"`
}

type AuthorizationRuleData struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"`
	ViewerID  string    `json:"viewerId"`
	Resource  string    `json:"resource"`
	Effect    string    `json:"effect"`
	CreatedAt time.Time `json:"createdAt"`
}

type AccountSettingsData struct {
	IsPrivate bool         `json:"isPrivate"`
	Devices   []DeviceData `json:"devices"`
}

type DeviceData struct {
	Token     string    `json:"token"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"createdAt"`
}

type InvitesData struct {
	SentInvites     []InviteData `json:"sentInvites"`
	ReceivedInvites []InviteData `json:"receivedInvites"`
}

type InviteData struct {
	ID             string    `json:"id"`
	OtherUserID    string    `json:"otherUserId"`
	OtherUsername  string    `json:"otherUsername"`
	OtherFirstName string    `json:"otherFirstName"`
	OtherLastName  string    `json:"otherLastName"`
	OtherImage     string    `json:"otherImage"`
	Status         string    `json:"status"`
	Note           *string   `json:"note,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
