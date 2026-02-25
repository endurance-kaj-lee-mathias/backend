package domain

import "github.com/gofrs/uuid"

type SampleId struct {
	uuid.UUID
}

func NewSampleId() (SampleId, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return SampleId{}, err
	}
	return SampleId{UUID: id}, nil
}
