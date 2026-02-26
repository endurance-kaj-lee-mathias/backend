package domain

import "github.com/gofrs/uuid"

type UserId struct {
	uuid.UUID
}

func NewUserId(value string) (UserId, error) {
	id, err := uuid.FromString(value)
	if err != nil {
		return UserId{}, err
	}
	return UserId{UUID: id}, nil
}
