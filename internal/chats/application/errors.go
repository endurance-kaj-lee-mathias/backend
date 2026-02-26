package application

import "errors"

var (
	NoSupportRelationship = errors.New("no support relationship exists between these users")
)
