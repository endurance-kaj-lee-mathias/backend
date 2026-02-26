package infrastructure

import "errors"

var (
	ConversationNotFound = errors.New("conversation not found")
	ParticipantNotFound  = errors.New("participant not found")
	UserNotFound         = errors.New("user not found")
)
