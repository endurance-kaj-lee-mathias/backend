package domain

import "errors"

var SelfSupportNotAllowed = errors.New("a user cannot support themselves")

var SelfInvite = errors.New("a user cannot invite themselves")
var DuplicatePendingInvite = errors.New("a pending invite already exists")
var AlreadyConnected = errors.New("a support relationship already exists")
var NotReceiver = errors.New("only the receiver can accept or decline an invite")
var InviteNotFound = errors.New("invite not found")
var NoteTooLong = errors.New("note must not exceed 300 characters")
