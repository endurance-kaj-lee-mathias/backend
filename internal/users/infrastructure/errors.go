package infrastructure

import "errors"

var NotFound = errors.New("user was not found")
var AddressNotFound = errors.New("address was not found")
