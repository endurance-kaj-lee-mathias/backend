package transport

import "errors"

var (
	Unauthorized = errors.New("unauthorized")
	InvalidId    = errors.New("invalid user id")
	ExportFailed = errors.New("failed to export user data")
)
