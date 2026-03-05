package models

import "errors"

var InvalidViewerID = errors.New("viewerId is required")
var InvalidResource = errors.New("invalid resource type")
var InvalidEffect = errors.New("effect must be 'allow' or 'deny'")
