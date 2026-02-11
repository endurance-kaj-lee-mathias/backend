package transport

import "errors"

var InvalidId = errors.New("id is invalid")
var InvalidQuery = errors.New("query is invalid")
var NotFound = errors.New("message was not found")
