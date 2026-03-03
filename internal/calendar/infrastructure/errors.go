package infrastructure

import "errors"

var SlotNotFound = errors.New("slot not found")
var AppointmentNotFound = errors.New("appointment not found")
var SlotOverlapDB = errors.New("slot overlaps with an existing slot")
var SlotAlreadyBookedDB = errors.New("slot is already booked")
