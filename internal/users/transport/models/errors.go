package models

import "errors"

var InvalidEmail = errors.New("email is invalid")
var InvalidStreet = errors.New("street is required")
var InvalidHouseNumber = errors.New("house number is required")
var InvalidPostalCode = errors.New("postal code is required")
var InvalidCity = errors.New("city is required")
var InvalidCountry = errors.New("country is required")
