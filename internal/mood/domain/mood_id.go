package domain

import "github.com/gofrs/uuid"

type MoodId struct {
	uuid.UUID
}

func NewMoodId() (MoodId, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return MoodId{}, err
	}
	return MoodId{UUID: id}, nil
}
