package models

import (
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserModel struct {
	ID          uuid.UUID     `json:"id"`
	Email       string        `json:"email"`
	Username    string        `json:"username"`
	FirstName   string        `json:"firstName"`
	LastName    string        `json:"lastName"`
	PhoneNumber *string       `json:"phoneNumber"`
	Roles       []domain.Role `json:"roles"`
}

func ToModel(usr domain.User) UserModel {
	return UserModel{
		ID:          usr.ID.UUID,
		Email:       usr.Email,
		Username:    usr.Username,
		FirstName:   usr.FirstName,
		LastName:    usr.LastName,
		PhoneNumber: usr.PhoneNumber,
		Roles:       usr.Roles,
	}
}
