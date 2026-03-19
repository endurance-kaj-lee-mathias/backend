package models

import "errors"

var InvalidStartTime = errors.New("startTime is required")
var InvalidEndTime = errors.New("endTime is required")
var InvalidFromParam = errors.New("from query parameter is required and must be a valid RFC3339 timestamp")
var InvalidToParam = errors.New("to query parameter is required and must be a valid RFC3339 timestamp")
var InvalidProviderIdParam = errors.New("providerId query parameter must be a valid UUID")
var InvalidDayParam = errors.New("day path parameter is required and must be a valid date in YYYY-MM-DD format")
