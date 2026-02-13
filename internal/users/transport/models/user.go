package models

import (
	"github.com/google/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserModel struct {
	ID    uuid.UUID     `json:"id"`
	Email string        `json:"email"`
	Roles []domain.Role `json:"roles"`
}

func ToModel(u domain.User) UserModel {
	return UserModel{
		ID:    u.ID,
		Email: u.Email,
		Roles: u.Roles,
	}
}
