package domain

import (
	"fmt"
	"time"
)

type RiskLevel string

const (
	RiskLevelNormal RiskLevel = "normal"
	RiskLevelHigh   RiskLevel = "high"
)

type User struct {
	ID           UserId    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	PhoneNumber  *string   `json:"phoneNumber"`
	Roles        []Role    `json:"roles"`
	About        string    `json:"about"`
	Introduction string    `json:"introduction"`
	Image        string    `json:"image"`
	RiskLevel    RiskLevel `json:"riskLevel"`
	IsPrivate    bool      `json:"isPrivate"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewUser(id UserId, email string, username string, firstName string, lastName string, roles []Role) User {
	now := time.Now().UTC()

	return User{
		ID:           id,
		Email:        email,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		Roles:        roles,
		About:        fmt.Sprintf("Hey, I'm %s. I'm here to find support and connect with others who understand.", firstName),
		Introduction: fmt.Sprintf("Hi, I'm %s! I'm new to Endurance and looking forward to connecting.", firstName),
		RiskLevel:    RiskLevelNormal,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
