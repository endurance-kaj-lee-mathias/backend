package models

import "errors"

var InvalidMoodScore = errors.New("mood_score must be between 0 and 10")
var InvalidDate = errors.New("date is required")
var NotesTooLong = errors.New("notes must not exceed 500 characters")
