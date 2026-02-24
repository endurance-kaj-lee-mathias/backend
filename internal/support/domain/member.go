package domain

import (
	"time"
)

type Member struct {
	ID        MemberId  `json:"id"`
	Veteran   VeteranId `json:"veteran"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewMember(id MemberId, veteran VeteranId, email string, username string, firstName string, lastName string, createdAt time.Time, updatedAt time.Time) Member {
	return Member{
		ID:        id,
		Veteran:   veteran,
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
