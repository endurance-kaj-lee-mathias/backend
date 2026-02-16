package domain

import "errors"

var ErrVeteranMustHaveVeteranRole = errors.New("only users with VETERAN role can receive support; THERAPIST and SUPPORT cannot be supported")
var ErrSupporterMustBeAbleToSupport = errors.New("user must have VETERAN, THERAPIST or SUPPORT role to provide support")
var ErrSelfSupportNotAllowed = errors.New("a user cannot support themselves")
