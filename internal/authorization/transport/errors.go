package transport

import "errors"

var Unauthorized = errors.New("unauthorized")
var InvalidId = errors.New("id is invalid")
