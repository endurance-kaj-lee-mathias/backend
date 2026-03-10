package application

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/sync/errgroup"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) ExportUserData(ctx context.Context, userID uuid.UUID) (domain.UserDataExport, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	userEnt, addrEnt, err := s.repo.GetUserWithAddress(ctx, userID)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	userKey, err := s.enc.DecryptUserKey(userEnt.EncryptedUserKey)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	profile, err := s.decryptProfile(userEnt, userKey)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	var address *domain.AddressData
	if addrEnt != nil {
		a, err := s.decryptAddress(*addrEnt, userKey)
		if err != nil {
			return domain.UserDataExport{}, err
		}
		address = &a
	}

	g, gCtx := errgroup.WithContext(ctx)

	var stressSamples []entities.StressSampleExportEntity
	var stressScores []entities.StressScoreExportEntity
	var moodEntries []entities.MoodEntryExportEntity
	var messageEnts []entities.MessageExportEntity
	var appointmentEnts []entities.AppointmentExportEntity
	var slotEnts []entities.SlotExportEntity
	var supporters []entities.SupportMemberExportEntity
	var supportedVeterans []entities.SupportMemberExportEntity
	var rulesOwned []entities.AuthorizationRuleExportEntity
	var rulesViewer []entities.AuthorizationRuleExportEntity
	var devices []entities.DeviceExportEntity
	var sentInvites []entities.InviteExportEntity
	var receivedInvites []entities.InviteExportEntity

	g.Go(func() error {
		var err error
		stressSamples, err = s.repo.GetStressSamples(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		stressScores, err = s.repo.GetStressScores(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		moodEntries, err = s.repo.GetMoodEntries(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		messageEnts, err = s.repo.GetConversationsAndMessages(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		appointmentEnts, err = s.repo.GetAppointmentsWithSlots(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		slotEnts, err = s.repo.GetSlotsAsProvider(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		supporters, err = s.repo.GetSupporters(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		supportedVeterans, err = s.repo.GetSupportedVeterans(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		rulesOwned, err = s.repo.GetAuthorizationRulesAsOwner(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		rulesViewer, err = s.repo.GetAuthorizationRulesAsViewer(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		devices, err = s.repo.GetDevices(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		sentInvites, err = s.repo.GetSentInvites(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		receivedInvites, err = s.repo.GetReceivedInvites(gCtx, userID)
		return err
	})

	if err := g.Wait(); err != nil {
		return domain.UserDataExport{}, err
	}

	decryptedSamples, err := s.decryptStressSamples(stressSamples, userKey)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	decryptedScores := s.mapStressScores(stressScores)

	decryptedMood, err := s.decryptMoodEntries(moodEntries, userKey)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	decryptedMessages, err := s.decryptMessages(messageEnts, userKey)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	decryptedAppointments, err := s.decryptAppointments(appointmentEnts)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	mappedSlots := s.mapSlots(slotEnts)

	decryptedSupporters, err := s.decryptSupportMembers(supporters)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	decryptedVeterans, err := s.decryptSupportMembers(supportedVeterans)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	mappedRulesOwned := s.mapAuthorizationRules(rulesOwned)
	mappedRulesViewer := s.mapAuthorizationRules(rulesViewer)

	mappedDevices := s.mapDevices(devices)

	decryptedSent, err := s.decryptInvites(sentInvites)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	decryptedReceived, err := s.decryptInvites(receivedInvites)
	if err != nil {
		return domain.UserDataExport{}, err
	}

	return domain.UserDataExport{
		ExportedAt: time.Now().UTC(),
		UserID:     userID.String(),
		Data: domain.ExportData{
			Profile: profile,
			Address: address,
			HealthData: domain.HealthData{
				StressSamples: decryptedSamples,
				StressScores:  decryptedScores,
				MoodEntries:   decryptedMood,
			},
			Messages: decryptedMessages,
			Calendar: domain.CalendarData{
				Appointments: decryptedAppointments,
				Slots:        mappedSlots,
			},
			SupportNetwork: domain.SupportNetworkData{
				Supporters:        decryptedSupporters,
				SupportedVeterans: decryptedVeterans,
			},
			DataSharing: domain.DataSharingData{
				SharedByMe:   mappedRulesOwned,
				SharedWithMe: mappedRulesViewer,
			},
			AccountSettings: domain.AccountSettingsData{
				IsPrivate: userEnt.IsPrivate,
				Devices:   mappedDevices,
			},
			Invites: domain.InvitesData{
				SentInvites:     decryptedSent,
				ReceivedInvites: decryptedReceived,
			},
		},
	}, nil
}
