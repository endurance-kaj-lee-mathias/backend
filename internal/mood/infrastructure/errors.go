package infrastructure

import "errors"

var UserNotFound = errors.New("user not found")
var MoodEntryNotFound = errors.New("mood entry not found")
