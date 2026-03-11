package transport

import "errors"

var InvalidId = errors.New("id is invalid")
var NotFound = errors.New("user was not found")
var Unauthorized = errors.New("unauthorized")
var AddressNotFound = errors.New("address was not found")
var UsernameAlreadyExists = errors.New("username already exists")
var DeviceNotFound = errors.New("device was not found")
var ClientNotEligible = errors.New("this client is not eligible for automatic role assignment")
