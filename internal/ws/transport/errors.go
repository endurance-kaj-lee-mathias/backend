package transport

import "errors"

var Unauthorized = errors.New("unauthorized")
var Forbidden = errors.New("forbidden")
var InvalidConversationID = errors.New("invalid conversation id")
