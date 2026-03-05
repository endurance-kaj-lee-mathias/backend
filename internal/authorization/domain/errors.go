package domain

import "errors"

var RuleNotFound = errors.New("authorization rule not found")
var NotOwner = errors.New("only the owner can manage authorization rules")
var InvalidResource = errors.New("invalid resource type")
var InvalidEffect = errors.New("invalid effect")
var SelfRule = errors.New("cannot create authorization rules for yourself")
