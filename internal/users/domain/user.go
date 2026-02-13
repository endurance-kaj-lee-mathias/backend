package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Roles     []Role    `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(email string, roles []Role) *User {
	now := time.Now().UTC()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Roles:     roles,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
