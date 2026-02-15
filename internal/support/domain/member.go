package domain

import (
	"time"
)

type Member struct {
	ID        MemberId  `json:"id"`
	Veteran   VeteranId `json:"veteran"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewMember(id MemberId, veteran VeteranId, email string) Member {
	now := time.Now().UTC()

	return Member{
		ID:        id,
		Veteran:   veteran,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
