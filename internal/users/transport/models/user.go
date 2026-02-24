package models

import (
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserModel struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	About        string    `json:"about"`
	Introduction string    `json:"introduction"`
	Image        string    `json:"image"`
}

func ToModel(usr domain.User) UserModel {
	return UserModel{
		ID:           usr.ID.UUID,
		Name:         usr.FirstName + " " + usr.LastName,
		Username:     usr.Username,
		About:        usr.About,
		Introduction: usr.Introduction,
		Image:        usr.Image,
	}
}
