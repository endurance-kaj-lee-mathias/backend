package domain

import "github.com/gofrs/uuid"

type AppointmentId struct {
	uuid.UUID
}

func NewAppointmentId() (AppointmentId, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return AppointmentId{}, err
	}
	return AppointmentId{UUID: id}, nil
}

func ParseAppointmentId(value string) (AppointmentId, error) {
	id, err := uuid.FromString(value)
	if err != nil {
		return AppointmentId{}, err
	}
	return AppointmentId{UUID: id}, nil
}
