package models

import "errors"

var InvalidEmail = errors.New("email is invalid")
var InvalidStreet = errors.New("street is required")
var InvalidLocality = errors.New("locality is required")
var InvalidPostalCode = errors.New("postal code is required")
var InvalidCountry = errors.New("country is required")
var InvalidIntroductionTooLong = errors.New("introduction must not exceed 500 characters")
var InvalidAboutTooLong = errors.New("about must not exceed 500 characters")
var InvalidImageEmpty = errors.New("image must not be empty")
var InvalidFirstNameEmpty = errors.New("first name must not be empty")
var InvalidLastNameEmpty = errors.New("last name must not be empty")
