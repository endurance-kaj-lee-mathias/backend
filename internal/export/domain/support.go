package domain

import (
	"time"
)

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
