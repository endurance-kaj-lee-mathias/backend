package domain

import (
	"github.com/gofrs/uuid"
)

type UserId struct {
	uuid.UUID
}

func NewId() UserId {
	return UserId{UUID: uuid.Must(uuid.NewV4())}
}

func ParseId(value string) (UserId, error) {
	id, err := uuid.FromString(value)

	if err != nil {
		return UserId{}, err
	}

	return UserId{UUID: id}, nil
}
