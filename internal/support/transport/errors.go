package transport

import "errors"

var InvalidId = errors.New("id is invalid")
var Unauthorized = errors.New("unauthorized")
var NotFound = errors.New("member was not found")
