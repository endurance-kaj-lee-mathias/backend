package transport

import "errors"

var Unauthorized = errors.New("unauthorized")
var InvalidId = errors.New("user_id is invalid")
var UserNotFound = errors.New("user was not found")
var UserIdMismatch = errors.New("user_id does not match authenticated user")
var ScoreNotFound = errors.New("no stress score found for this user")
