package infrastructure

import "errors"

var UserNotFound = errors.New("user not found")
var ScoreNotFound = errors.New("stress score not found")
var SampleNotFound = errors.New("stress sample not found")
var AlgoServiceUnavailable = errors.New("algorithm service is unavailable")
var AlgoServiceError = errors.New("algorithm service returned an error")
