package domain

import (
	"time"
)

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
