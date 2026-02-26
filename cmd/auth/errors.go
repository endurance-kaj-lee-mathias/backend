package auth

import "errors"

var ClaimsInvalid = errors.New("claims could not be parsed")
var HeaderInvalid = errors.New("authentication header is invalid")
var TokenInvalid = errors.New("token is invalid")
var MissingHeader = errors.New("authentication header is missing")
var IssuerInvalid = errors.New("token issuer is not trusted")
var MissingRole = errors.New("role is missing")
