package request

import "errors"

var EmptyBody = errors.New("request body is empty")
var InvalidJSON = errors.New("invalid JSON format")
var BodyTooLarge = errors.New("request body too large")
var ContentType = errors.New("content-type must be application/json")
