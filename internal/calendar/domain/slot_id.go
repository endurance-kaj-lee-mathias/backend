package domain

import "github.com/gofrs/uuid"

type SlotId struct {
	uuid.UUID
}

func NewSlotId() (SlotId, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return SlotId{}, err
	}
	return SlotId{UUID: id}, nil
}

func ParseSlotId(value string) (SlotId, error) {
	id, err := uuid.FromString(value)
	if err != nil {
		return SlotId{}, err
	}
	return SlotId{UUID: id}, nil
}
