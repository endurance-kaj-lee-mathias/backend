package transport

import "errors"

var InvalidId = errors.New("id is invalid")
var NotFound = errors.New("user was not found")
var Unauthorized = errors.New("unauthorized")
