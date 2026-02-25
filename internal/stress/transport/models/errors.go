package models

import "errors"

var InvalidUserID = errors.New("user_id is required")
var InvalidTimestamp = errors.New("timestamp_utc is required")
var InvalidWindowMinutes = errors.New("window_minutes must be greater than 0")
var InvalidMeanHR = errors.New("mean_hr must be greater than 0")
var InvalidRMSSD = errors.New("rmssd_ms must be greater than 0")
