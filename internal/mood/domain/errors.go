package domain

import "errors"

var InvalidScore = errors.New("mood_score must be between 1 and 10")
var NotesTooLong = errors.New("notes must not exceed 500 characters")
