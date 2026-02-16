package domain

import "errors"

var VeteranMustHaveVeteranRole = errors.New("only users with VETERAN role can receive support; THERAPIST and SUPPORT cannot be supported")
var SupporterMustBeAbleToSupport = errors.New("user must have VETERAN, THERAPIST or SUPPORT role to provide support")
var SelfSupportNotAllowed = errors.New("a user cannot support themselves")
