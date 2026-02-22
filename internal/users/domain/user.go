package domain

import (
	"time"
)

type User struct {
	ID          UserId    `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber *string   `json:"phoneNumber"`
	Roles       []Role    `json:"roles"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewUser(id UserId, email string, firstName string, lastName string, roles []Role) User {
	now := time.Now().UTC()

	return User{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Roles:     roles,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
