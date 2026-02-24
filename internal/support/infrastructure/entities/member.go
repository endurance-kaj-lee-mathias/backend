package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
)

type MemberEntity struct {
	ID               uuid.UUID `db:"id"`
	Veteran          uuid.UUID `db:"veteran"`
	EncryptedEmail   []byte    `db:"encrypted_email"`
	EncryptedFirst   []byte    `db:"encrypted_first_name"`
	EncryptedLast    []byte    `db:"encrypted_last_name"`
	EncryptedUserKey []byte    `db:"encrypted_user_key"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func FromEntity(ent MemberEntity, enc encryption.Service) (domain.Member, error) {
	userKey, err := enc.DecryptUserKey(ent.EncryptedUserKey)
	if err != nil {
		return domain.Member{}, err
	}

	emailBytes, err := enc.Decrypt(ent.EncryptedEmail, userKey)
	if err != nil {
		return domain.Member{}, err
	}

	firstNameBytes, err := enc.Decrypt(ent.EncryptedFirst, userKey)
	if err != nil {
		return domain.Member{}, err
	}

	lastNameBytes, err := enc.Decrypt(ent.EncryptedLast, userKey)
	if err != nil {
		return domain.Member{}, err
	}

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
		string(emailBytes),
		string(firstNameBytes),
		string(lastNameBytes),
		ent.CreatedAt,
		ent.UpdatedAt,
	), nil
}

func FromEntities(ents []MemberEntity, enc encryption.Service) ([]domain.Member, error) {
	out := make([]domain.Member, 0, len(ents))

	for _, ent := range ents {
		member, err := FromEntity(ent, enc)
		if err != nil {
			return nil, err
		}

		out = append(out, member)
	}

	return out, nil
}
