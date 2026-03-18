package domain

import "errors"

var MissingUserProfile = errors.New("user profile is required")
var NotVeteran = errors.New("target user is not a veteran")
