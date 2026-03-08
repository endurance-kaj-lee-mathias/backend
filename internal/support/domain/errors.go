package domain

import "errors"

var VeteranMustHaveVeteranRole = errors.New("only users with VETERAN role can receive support; THERAPIST and SUPPORT cannot be supported")
var SupporterMustBeAbleToSupport = errors.New("user must have VETERAN, THERAPIST or SUPPORT role to provide support")
var SelfSupportNotAllowed = errors.New("a user cannot support themselves")

var SelfInvite = errors.New("a user cannot invite themselves")
var DuplicatePendingInvite = errors.New("a pending invite already exists")
var AlreadyAccepted = errors.New("an accepted relationship already exists")
var NotReceiver = errors.New("only the receiver can accept or decline an invite")
var InviteNotFound = errors.New("invite not found")
var NoteTooLong = errors.New("note must not exceed 300 characters")
