package transport

import "errors"

var (
	Unauthorized          = errors.New("unauthorized")
	InvalidConversationID = errors.New("invalid conversation ID")
	InvalidRequestBody    = errors.New("invalid request body")
	Forbidden             = errors.New("forbidden")
)
