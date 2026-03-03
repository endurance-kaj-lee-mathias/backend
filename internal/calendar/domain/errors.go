package domain

import "errors"

var InvalidTimeRange = errors.New("end time must be after start time")
var SlotInPast = errors.New("cannot create a slot in the past")
var SlotOverlap = errors.New("slot overlaps with an existing slot")
var SlotAlreadyBooked = errors.New("slot is already booked")
var SlotNotFound = errors.New("slot not found")
var AppointmentNotFound = errors.New("appointment not found")
var CannotDeleteBookedSlot = errors.New("cannot delete a booked slot")
var NormalCannotBookUrgent = errors.New("normal booking cannot use an urgent slot")
var UrgentRequiresUrgentSlot = errors.New("urgent booking requires an urgent slot")
var OnlyProviderCanManageSlots = errors.New("only therapists and support can manage slots")
var OnlyVeteranCanBook = errors.New("only veterans can book slots")
var NotSlotOwner = errors.New("you can only manage your own slots")
var NotAppointmentParticipant = errors.New("only the booking veteran or slot provider can cancel")
var InsufficientUrgentSlots = errors.New("you must create enough urgent slots for this day before adding non-urgent slots")
var SlotInPastCannotBook = errors.New("cannot book a slot in the past")
