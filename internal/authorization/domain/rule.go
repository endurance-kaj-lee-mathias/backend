package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Rule struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	ViewerID  uuid.UUID
	Resource  ResourceType
	Effect    PolicyEffect
	CreatedAt time.Time
}

func NewRule(ownerID uuid.UUID, viewerID uuid.UUID, resource ResourceType, effect PolicyEffect) (Rule, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Rule{}, err
	}

	return Rule{
		ID:        id,
		OwnerID:   ownerID,
		ViewerID:  viewerID,
		Resource:  resource,
		Effect:    effect,
		CreatedAt: time.Now().UTC(),
	}, nil
}
