package domain

import "github.com/gofrs/uuid"

type AddressId struct {
	uuid.UUID
}

func NewAddressId() (AddressId, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return AddressId{}, err
	}
	return AddressId{UUID: id}, nil
}

func ParseAddressId(value string) (AddressId, error) {
	id, err := uuid.FromString(value)
	if err != nil {
		return AddressId{}, err
	}
	return AddressId{UUID: id}, nil
}
