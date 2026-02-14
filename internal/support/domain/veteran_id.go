package domain

import (
	"github.com/gofrs/uuid"
)

type VeteranId struct {
	uuid.UUID
}

func NewVeteranId() VeteranId {
	return VeteranId{UUID: uuid.Must(uuid.NewV4())}
}

func ParseVeteranId(value string) (VeteranId, error) {
	id, err := uuid.FromString(value)

	if err != nil {
		return VeteranId{}, err
	}

	return VeteranId{UUID: id}, nil
}
