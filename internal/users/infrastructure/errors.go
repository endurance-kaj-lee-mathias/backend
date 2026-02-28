package infrastructure

import "errors"

var NotFound = errors.New("user was not found")
var AddressNotFound = errors.New("address was not found")
var UsernameAlreadyExists = errors.New("username already exists")
var DeviceNotFound = errors.New("device was not found")
