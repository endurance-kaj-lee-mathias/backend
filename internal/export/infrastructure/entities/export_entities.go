package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type UserExportEntity struct {
	ID                    uuid.UUID
	EncryptedEmail        []byte
	EncryptedUsername     []byte
	EncryptedFirstName    []byte
	EncryptedLastName     []byte
	EncryptedPhoneNumber  []byte
	EncryptedRoles        []byte
	EncryptedAbout        []byte
	EncryptedIntroduction []byte
	Image                 string
	IsPrivate             bool
	EncryptedUserKey      []byte
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type AddressExportEntity struct {
	ID                  uuid.UUID
	EncryptedStreet     []byte
	EncryptedLocality   []byte
	EncryptedRegion     []byte
	EncryptedPostalCode []byte
	EncryptedCountry    []byte
	CreatedAt           time.Time
}

type StressSampleExportEntity struct {
	ID                      uuid.UUID
	TimestampUTC            time.Time
	WindowMinutes           int
	EncryptedMeanHR         []byte
	EncryptedRMSSDms        []byte
	EncryptedRestingHR      []byte
	EncryptedSteps          []byte
	EncryptedSleepDebtHours []byte
	CreatedAt               time.Time
}

type StressScoreExportEntity struct {
	ID           uuid.UUID
	Score        float64
	Category     string
	ModelVersion string
	ComputedAt   time.Time
}

type MoodEntryExportEntity struct {
	ID             uuid.UUID
	Date           time.Time
	MoodScore      int
	EncryptedNotes []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MessageExportEntity struct {
	ID                                 uuid.UUID
	ConversationID                     uuid.UUID
	SenderID                           uuid.UUID
	EncryptedContent                   []byte
	CreatedAt                          time.Time
	EncryptedConversationKey           []byte
	OtherParticipantID                 uuid.UUID
	OtherParticipantEncryptedUsername  []byte
	OtherParticipantEncryptedFirstName []byte
	OtherParticipantEncryptedLastName  []byte
	OtherParticipantEncryptedUserKey   []byte
}

type AppointmentExportEntity struct {
	ID                        uuid.UUID
	SlotID                    uuid.UUID
	VeteranID                 uuid.UUID
	ProviderID                uuid.UUID
	Status                    string
	StartTime                 time.Time
	EndTime                   time.Time
	IsUrgent                  bool
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	ProviderEncryptedUsername []byte
	ProviderEncryptedUserKey  []byte
}

type SlotExportEntity struct {
	ID        uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	IsUrgent  bool
	IsBooked  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SupportMemberExportEntity struct {
	ID                uuid.UUID
	EncryptedEmail    []byte
	EncryptedUsername []byte
	EncryptedFirst    []byte
	EncryptedLast     []byte
	Image             string
	EncryptedUserKey  []byte
	CreatedAt         time.Time
}

type AuthorizationRuleExportEntity struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	ViewerID  uuid.UUID
	Resource  string
	Effect    string
	CreatedAt time.Time
}

type DeviceExportEntity struct {
	DeviceToken string
	Platform    string
	CreatedAt   time.Time
}

type InviteExportEntity struct {
	ID                      uuid.UUID
	OtherUserID             uuid.UUID
	OtherEncryptedUsername  []byte
	OtherEncryptedFirstName []byte
	OtherEncryptedLastName  []byte
	OtherEncryptedUserKey   []byte
	OtherImage              *string
	Status                  string
	Note                    *string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
