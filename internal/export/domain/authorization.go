package domain

import (
	"time"
)

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
