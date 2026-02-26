package models

import "errors"

var (
	InvalidParticipantID = errors.New("invalid participant ID")
	EmptyContent         = errors.New("message content cannot be empty")
)
