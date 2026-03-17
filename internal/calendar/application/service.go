package application

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure/entities"
)

func (s *service) CreateSlot(ctx context.Context, providerID uuid.UUID, roles []string, startTime, endTime time.Time, isUrgent bool, isRecurring bool) (domain.Slot, error) {
	if !domain.HasProviderRole(roles) {
		return domain.Slot{}, domain.OnlyProviderCanManageSlots
	}

	if err := domain.ValidateSlotTimes(startTime, endTime); err != nil {
		return domain.Slot{}, err
	}

	if !isUrgent && domain.HasTherapistRole(roles) {
		urgentMinutes, err := s.repo.GetUrgentSlotMinutesForDate(ctx, providerID, startTime)
		if err != nil {
			return domain.Slot{}, err
		}
		if urgentMinutes < s.minUrgentMinutes {
			return domain.Slot{}, domain.InsufficientUrgentSlots
		}
	}

	var seriesID *uuid.UUID
	if isRecurring {
		id, err := uuid.NewV4()
		if err != nil {
			return domain.Slot{}, err
		}
		seriesID = &id
	}

	duration := endTime.Sub(startTime)
	weeks := 1
	if isRecurring {
		weeks = 8
	}

	var firstSlot domain.Slot
	now := time.Now().UTC()
	for i := 0; i < weeks; i++ {
		slotID, err := domain.NewSlotId()
		if err != nil {
			return domain.Slot{}, err
		}
		offset := time.Duration(i) * 7 * 24 * time.Hour
		slot := domain.Slot{
			ID:         slotID,
			ProviderID: providerID,
			StartTime:  startTime.Add(offset),
			EndTime:    startTime.Add(offset).Add(duration),
			IsUrgent:   isUrgent,
			IsBooked:   false,
			SeriesID:   seriesID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		ent := entities.SlotToEntity(slot)
		if err := s.repo.CreateSlot(ctx, ent); err != nil {
			if errors.Is(err, infrastructure.SlotOverlapDB) {
				return domain.Slot{}, domain.SlotOverlap
			}
			return domain.Slot{}, err
		}
		if i == 0 {
			firstSlot = slot
		}
	}

	return firstSlot, nil
}

func (s *service) GetSlots(ctx context.Context, from, to time.Time, providerID *uuid.UUID) ([]domain.Slot, error) {
	ents, err := s.repo.GetSlotsByRange(ctx, from, to, providerID)
	if err != nil {
		return nil, err
	}
	return entities.SlotsFromEntities(ents), nil
}

func (s *service) DeleteSlot(ctx context.Context, userID uuid.UUID, roles []string, slotID domain.SlotId) error {
	if !domain.HasProviderRole(roles) {
		return domain.OnlyProviderCanManageSlots
	}

	ent, err := s.repo.GetSlotByID(ctx, slotID.UUID)
	if err != nil {
		if errors.Is(err, infrastructure.SlotNotFound) {
			return domain.SlotNotFound
		}
		return err
	}

	if ent.ProviderID != userID {
		return domain.NotSlotOwner
	}

	if ent.IsBooked {
		return domain.CannotDeleteBookedSlot
	}

	if err := s.repo.DeleteSlot(ctx, slotID.UUID); err != nil {
		if errors.Is(err, infrastructure.SlotNotFound) {
			return domain.SlotNotFound
		}
		return err
	}

	return nil
}

func (s *service) BookSlot(ctx context.Context, veteranID uuid.UUID, roles []string, slotID domain.SlotId, urgent bool) (domain.Appointment, error) {
	if !domain.HasVeteranRole(roles) {
		return domain.Appointment{}, domain.OnlyVeteranCanBook
	}

	ent, err := s.repo.GetSlotByID(ctx, slotID.UUID)
	if err != nil {
		if errors.Is(err, infrastructure.SlotNotFound) {
			return domain.Appointment{}, domain.SlotNotFound
		}
		return domain.Appointment{}, err
	}

	slot := entities.SlotFromEntity(ent)
	if err := domain.ValidateBooking(slot, urgent); err != nil {
		return domain.Appointment{}, err
	}

	now := time.Now().UTC()
	rowsAffected, err := s.repo.AtomicBookSlot(ctx, slotID.UUID, now)
	if err != nil {
		return domain.Appointment{}, err
	}
	if rowsAffected == 0 {
		return domain.Appointment{}, domain.SlotAlreadyBooked
	}

	appointmentID, err := domain.NewAppointmentId()
	if err != nil {
		return domain.Appointment{}, err
	}

	appointment := domain.Appointment{
		ID:        appointmentID,
		SlotID:    slotID.UUID,
		VeteranID: veteranID,
		Status:    domain.StatusBooked,
		CreatedAt: now,
		UpdatedAt: now,
	}

	appointmentEnt := entities.AppointmentToEntity(appointment)
	if err := s.repo.CreateAppointment(ctx, appointmentEnt); err != nil {
		return domain.Appointment{}, err
	}

	return appointment, nil
}

func (s *service) CancelAppointment(ctx context.Context, userID uuid.UUID, appointmentID domain.AppointmentId) error {
	ent, err := s.repo.GetAppointmentWithSlot(ctx, appointmentID.UUID)
	if err != nil {
		if errors.Is(err, infrastructure.AppointmentNotFound) {
			return domain.AppointmentNotFound
		}
		return err
	}

	if ent.VeteranID != userID && ent.SlotProviderID != userID {
		return domain.NotAppointmentParticipant
	}

	return s.repo.CancelAppointment(ctx, appointmentID.UUID, time.Now().UTC())
}

func (s *service) DeleteMySlots(ctx context.Context, providerID uuid.UUID) error {
	return s.repo.DeleteSlotsByProviderID(ctx, providerID)
}

func (s *service) DeleteSlotsBySeries(ctx context.Context, providerID uuid.UUID, seriesID uuid.UUID) error {
	return s.repo.DeleteFutureSlotsBySeries(ctx, seriesID, providerID)
}
