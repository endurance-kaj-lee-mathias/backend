package domain

import (
	"time"
)

type User struct {
	ID          UserId    `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber *string   `json:"phoneNumber"`
	Roles       []Role    `json:"roles"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewUser(id UserId, email string, username string, firstName string, lastName string, roles []Role) User {
	now := time.Now().UTC()

	return User{
		ID:        id,
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Roles:     roles,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
