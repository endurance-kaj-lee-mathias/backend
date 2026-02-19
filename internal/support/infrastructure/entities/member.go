package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
)

type MemberEntity struct {
	ID        uuid.UUID `db:"id"`
	Veteran   uuid.UUID `db:"veteran"`
	Email     string    `db:"email"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func FromEntity(ent MemberEntity) (domain.Member, error) {
	memberId, err := domain.ParseMemberId(ent.ID.String())

	if err != nil {
		return domain.Member{}, err
	}

	veteranId, err := domain.ParseVeteranId(ent.Veteran.String())

	if err != nil {
		return domain.Member{}, err
	}

	return domain.NewMember(
		memberId,
		veteranId,
		ent.Email,
		ent.FirstName,
		ent.LastName,
		ent.CreatedAt,
		ent.UpdatedAt,
	), nil
}

func FromEntities(ents []MemberEntity) []domain.Member {
	out := make([]domain.Member, 0, len(ents))

	for _, ent := range ents {
		ent, err := FromEntity(ent)

		if err != nil {
			continue
		}

		out = append(out, ent)
	}

	return out
}

func ToEntity(mem domain.Member) (MemberEntity, error) {
	return MemberEntity{
		ID:        mem.ID.UUID,
		Veteran:   mem.Veteran.UUID,
		Email:     mem.Email,
		FirstName: mem.FirstName,
		LastName:  mem.LastName,
		CreatedAt: mem.CreatedAt,
		UpdatedAt: mem.UpdatedAt,
	}, nil
}
