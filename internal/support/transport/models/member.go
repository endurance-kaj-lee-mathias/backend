package models

import (
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
)

type MemberModel struct {
	ID        uuid.UUID `json:"id"`
	Veteran   uuid.UUID `json:"veteran"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
}

func ToModel(mem domain.Member) MemberModel {
	return MemberModel{
		ID:        mem.ID.UUID,
		Veteran:   mem.Veteran.UUID,
		Email:     mem.Email,
		Username:  mem.Username,
		FirstName: mem.FirstName,
		LastName:  mem.LastName,
	}
}

func ToModels(mems []domain.Member) []MemberModel {
	out := make([]MemberModel, 0, len(mems))

	for _, reg := range mems {
		out = append(out, ToModel(reg))
	}

	return out
}
