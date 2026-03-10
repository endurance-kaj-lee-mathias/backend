package transport

import "errors"

var Unauthorized = errors.New("unauthorized")
var InvalidId = errors.New("user_id is invalid")
var InvalidEntryId = errors.New("entry_id is invalid")
var UserNotFound = errors.New("user was not found")
var InvalidDate = errors.New("date format must be YYYY-MM-DD")
var MoodEntryNotFound = errors.New("mood entry not found")
var Forbidden = errors.New("you are not allowed to modify this entry")
