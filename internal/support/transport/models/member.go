package models

import (
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
)

type MemberModel struct {
	ID      uuid.UUID `json:"id"`
	Veteran uuid.UUID `json:"veteran"`
	Email   string    `json:"email"`
}

func ToModel(mem domain.Member) MemberModel {
	return MemberModel{
		ID:      mem.ID.UUID,
		Veteran: mem.Veteran.UUID,
		Email:   mem.Email,
	}
}

func ToModels(mems []domain.Member) []MemberModel {
	out := make([]MemberModel, 0, len(mems))

	for _, reg := range mems {
		out = append(out, ToModel(reg))
	}

	return out
}
