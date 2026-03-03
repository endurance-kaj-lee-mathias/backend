package domain

import "github.com/gofrs/uuid"

type ScoreId struct {
	uuid.UUID
}

func NewScoreId() (ScoreId, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return ScoreId{}, err
	}
	return ScoreId{UUID: id}, nil
}
