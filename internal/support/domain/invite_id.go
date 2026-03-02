package domain

import "github.com/gofrs/uuid"

type InviteId struct {
	uuid.UUID
}

func NewInviteId() InviteId {
	return InviteId{UUID: uuid.Must(uuid.NewV4())}
}

func ParseInviteId(value string) (InviteId, error) {
	id, err := uuid.FromString(value)
	if err != nil {
		return InviteId{}, err
	}
	return InviteId{UUID: id}, nil
}
