package domain

import (
	"github.com/gofrs/uuid"
)

type MemberId struct {
	uuid.UUID
}

func NewMemberId() MemberId {
	return MemberId{UUID: uuid.Must(uuid.NewV4())}
}

func ParseMemberId(value string) (MemberId, error) {
	id, err := uuid.FromString(value)

	if err != nil {
		return MemberId{}, err
	}

	return MemberId{UUID: id}, nil
}
