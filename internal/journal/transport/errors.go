package transport

import "errors"

var Unauthorized = errors.New("unauthorized")
var InvalidId = errors.New("id is invalid")
var InvalidUsername = errors.New("username is invalid")
var VeteranNotFound = errors.New("veteran not found")
var Forbidden = errors.New("forbidden")
