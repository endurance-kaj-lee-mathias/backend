package transport

import "errors"

var Unauthorized = errors.New("unauthorized")
var InvalidId = errors.New("user_id is invalid")
var UserNotFound = errors.New("user was not found")
var InvalidDate = errors.New("date format must be YYYY-MM-DD")
