package models

import "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"

type AddUserModel struct {
	Email string        `json:"email"`
	Roles []domain.Role `json:"roles"`
}

func (mod *AddUserModel) Validate() error {
	if mod.Email == "" {
		return InvalidEmail
	}

	return nil
}
