package application

import (
	"context"
	"strconv"
	"strings"
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

func (s *service) decryptProfile(ent entities.UserExportEntity, userKey []byte) (domain.ProfileData, error) {
	email, err := s.enc.Decrypt(ent.EncryptedEmail, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	username, err := s.enc.Decrypt(ent.EncryptedUsername, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	firstName, err := s.enc.Decrypt(ent.EncryptedFirstName, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	lastName, err := s.enc.Decrypt(ent.EncryptedLastName, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	about, err := s.enc.Decrypt(ent.EncryptedAbout, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	introduction, err := s.enc.Decrypt(ent.EncryptedIntroduction, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	rolesStr, err := s.enc.Decrypt(ent.EncryptedRoles, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}
	roles := strings.Split(string(rolesStr), ",")
	if len(roles) == 1 && roles[0] == "" {
		roles = []string{}
	}

	var phoneNumber *string
	if len(ent.EncryptedPhoneNumber) > 0 {
		phone, err := s.enc.Decrypt(ent.EncryptedPhoneNumber, userKey)
		if err != nil {
			return domain.ProfileData{}, err
		}
		p := string(phone)
		phoneNumber = &p
	}

	return domain.ProfileData{
		Email:        string(email),
		Username:     string(username),
		FirstName:    string(firstName),
		LastName:     string(lastName),
		PhoneNumber:  phoneNumber,
		Roles:        roles,
		About:        string(about),
		Introduction: string(introduction),
		Image:        ent.Image,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,
	}, nil
}

func (s *service) decryptAddress(ent entities.AddressExportEntity, userKey []byte) (domain.AddressData, error) {
	street, err := s.enc.Decrypt(ent.EncryptedStreet, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	locality, err := s.enc.Decrypt(ent.EncryptedLocality, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	region, err := s.enc.Decrypt(ent.EncryptedRegion, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	postalCode, err := s.enc.Decrypt(ent.EncryptedPostalCode, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	country, err := s.enc.Decrypt(ent.EncryptedCountry, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	return domain.AddressData{
		Street:     string(street),
		Locality:   string(locality),
		Region:     string(region),
		PostalCode: string(postalCode),
		Country:    string(country),
		CreatedAt:  ent.CreatedAt,
	}, nil
}

func (s *service) decryptStressSamples(ents []entities.StressSampleExportEntity, userKey []byte) ([]domain.StressSampleData, error) {
	result := make([]domain.StressSampleData, 0, len(ents))

	for _, ent := range ents {
		meanHR, err := s.decryptFloat(ent.EncryptedMeanHR, userKey)
		if err != nil {
			return nil, err
		}

		rmssd, err := s.decryptFloat(ent.EncryptedRMSSDms, userKey)
		if err != nil {
			return nil, err
		}

		var restingHR *float64
		if len(ent.EncryptedRestingHR) > 0 {
			val, err := s.decryptFloat(ent.EncryptedRestingHR, userKey)
			if err != nil {
				return nil, err
			}
			restingHR = &val
		}

		var steps *int
		if len(ent.EncryptedSteps) > 0 {
			val, err := s.decryptInt(ent.EncryptedSteps, userKey)
			if err != nil {
				return nil, err
			}
			steps = &val
		}

		var sleepDebt *float64
		if len(ent.EncryptedSleepDebtHours) > 0 {
			val, err := s.decryptFloat(ent.EncryptedSleepDebtHours, userKey)
			if err != nil {
				return nil, err
			}
			sleepDebt = &val
		}

		result = append(result, domain.StressSampleData{
			ID:             ent.ID.String(),
			TimestampUTC:   ent.TimestampUTC,
			WindowMinutes:  ent.WindowMinutes,
			MeanHR:         meanHR,
			RMSSDms:        rmssd,
			RestingHR:      restingHR,
			Steps:          steps,
			SleepDebtHours: sleepDebt,
			CreatedAt:      ent.CreatedAt,
		})
	}

	return result, nil
}

func (s *service) mapStressScores(ents []entities.StressScoreExportEntity) []domain.StressScoreData {
	result := make([]domain.StressScoreData, 0, len(ents))

	for _, ent := range ents {
		result = append(result, domain.StressScoreData{
			ID:           ent.ID.String(),
			Score:        ent.Score,
			Category:     ent.Category,
			ModelVersion: ent.ModelVersion,
			ComputedAt:   ent.ComputedAt,
		})
	}

	return result
}

func (s *service) decryptMoodEntries(ents []entities.MoodEntryExportEntity, userKey []byte) ([]domain.MoodEntryData, error) {
	result := make([]domain.MoodEntryData, 0, len(ents))

	for _, ent := range ents {
		var notes *string
		if len(ent.EncryptedNotes) > 0 {
			decrypted, err := s.enc.Decrypt(ent.EncryptedNotes, userKey)
			if err != nil {
				return nil, err
			}
			n := string(decrypted)
			notes = &n
		}

		result = append(result, domain.MoodEntryData{
			ID:        ent.ID.String(),
			Date:      ent.Date,
			MoodScore: ent.MoodScore,
			Notes:     notes,
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
		})
	}

	return result, nil
}

func (s *service) decryptMessages(ents []entities.MessageExportEntity, userKey []byte) ([]domain.MessageData, error) {
	result := make([]domain.MessageData, 0, len(ents))
	convKeyCache := make(map[uuid.UUID][]byte)
	otherUserCache := make(map[uuid.UUID][3]string)

	for _, ent := range ents {
		convKey, ok := convKeyCache[ent.ConversationID]
		if !ok {
			decryptedConvKey, err := s.enc.Decrypt(ent.EncryptedConversationKey, userKey)
			if err != nil {
				return nil, err
			}
			convKey = decryptedConvKey
			convKeyCache[ent.ConversationID] = convKey
		}

		content, err := s.enc.Decrypt(ent.EncryptedContent, convKey)
		if err != nil {
			return nil, err
		}

		otherInfo, ok := otherUserCache[ent.OtherParticipantID]
		if !ok {
			otherUserKey, err := s.enc.DecryptUserKey(ent.OtherParticipantEncryptedUserKey)
			if err != nil {
				return nil, err
			}

			username, err := s.enc.Decrypt(ent.OtherParticipantEncryptedUsername, otherUserKey)
			if err != nil {
				return nil, err
			}

			firstName, err := s.enc.Decrypt(ent.OtherParticipantEncryptedFirstName, otherUserKey)
			if err != nil {
				return nil, err
			}

			lastName, err := s.enc.Decrypt(ent.OtherParticipantEncryptedLastName, otherUserKey)
			if err != nil {
				return nil, err
			}

			otherInfo = [3]string{string(username), string(firstName), string(lastName)}
			otherUserCache[ent.OtherParticipantID] = otherInfo
		}

		result = append(result, domain.MessageData{
			ID:                        ent.ID.String(),
			ConversationID:            ent.ConversationID.String(),
			SenderID:                  ent.SenderID.String(),
			Content:                   string(content),
			CreatedAt:                 ent.CreatedAt,
			OtherParticipantUsername:  otherInfo[0],
			OtherParticipantFirstName: otherInfo[1],
			OtherParticipantLastName:  otherInfo[2],
		})
	}

	return result, nil
}

func (s *service) decryptAppointments(ents []entities.AppointmentExportEntity) ([]domain.AppointmentData, error) {
	result := make([]domain.AppointmentData, 0, len(ents))
	providerCache := make(map[uuid.UUID]string)

	for _, ent := range ents {
		providerUsername, ok := providerCache[ent.ProviderID]
		if !ok {
			providerKey, err := s.enc.DecryptUserKey(ent.ProviderEncryptedUserKey)
			if err != nil {
				return nil, err
			}

			username, err := s.enc.Decrypt(ent.ProviderEncryptedUsername, providerKey)
			if err != nil {
				return nil, err
			}

			providerUsername = string(username)
			providerCache[ent.ProviderID] = providerUsername
		}

		result = append(result, domain.AppointmentData{
			ID:               ent.ID.String(),
			SlotID:           ent.SlotID.String(),
			VeteranID:        ent.VeteranID.String(),
			ProviderID:       ent.ProviderID.String(),
			ProviderUsername: providerUsername,
			Status:           ent.Status,
			StartTime:        ent.StartTime,
			EndTime:          ent.EndTime,
			IsUrgent:         ent.IsUrgent,
			CreatedAt:        ent.CreatedAt,
			UpdatedAt:        ent.UpdatedAt,
		})
	}

	return result, nil
}

func (s *service) mapSlots(ents []entities.SlotExportEntity) []domain.SlotData {
	result := make([]domain.SlotData, 0, len(ents))

	for _, ent := range ents {
		result = append(result, domain.SlotData{
			ID:        ent.ID.String(),
			StartTime: ent.StartTime,
			EndTime:   ent.EndTime,
			IsUrgent:  ent.IsUrgent,
			IsBooked:  ent.IsBooked,
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
		})
	}

	return result
}

func (s *service) decryptSupportMembers(ents []entities.SupportMemberExportEntity) ([]domain.SupportMemberData, error) {
	result := make([]domain.SupportMemberData, 0, len(ents))

	for _, ent := range ents {
		memberKey, err := s.enc.DecryptUserKey(ent.EncryptedUserKey)
		if err != nil {
			return nil, err
		}

		email, err := s.enc.Decrypt(ent.EncryptedEmail, memberKey)
		if err != nil {
			return nil, err
		}

		username, err := s.enc.Decrypt(ent.EncryptedUsername, memberKey)
		if err != nil {
			return nil, err
		}

		firstName, err := s.enc.Decrypt(ent.EncryptedFirst, memberKey)
		if err != nil {
			return nil, err
		}

		lastName, err := s.enc.Decrypt(ent.EncryptedLast, memberKey)
		if err != nil {
			return nil, err
		}

		result = append(result, domain.SupportMemberData{
			ID:        ent.ID.String(),
			Email:     string(email),
			Username:  string(username),
			FirstName: string(firstName),
			LastName:  string(lastName),
			Image:     ent.Image,
			CreatedAt: ent.CreatedAt,
		})
	}

	return result, nil
}

func (s *service) mapAuthorizationRules(ents []entities.AuthorizationRuleExportEntity) []domain.AuthorizationRuleData {
	result := make([]domain.AuthorizationRuleData, 0, len(ents))

	for _, ent := range ents {
		result = append(result, domain.AuthorizationRuleData{
			ID:        ent.ID.String(),
			OwnerID:   ent.OwnerID.String(),
			ViewerID:  ent.ViewerID.String(),
			Resource:  ent.Resource,
			Effect:    ent.Effect,
			CreatedAt: ent.CreatedAt,
		})
	}

	return result
}

func (s *service) mapDevices(ents []entities.DeviceExportEntity) []domain.DeviceData {
	result := make([]domain.DeviceData, 0, len(ents))

	for _, ent := range ents {
		result = append(result, domain.DeviceData{
			Token:     ent.DeviceToken,
			Platform:  ent.Platform,
			CreatedAt: ent.CreatedAt,
		})
	}

	return result
}

func (s *service) decryptInvites(ents []entities.InviteExportEntity) ([]domain.InviteData, error) {
	result := make([]domain.InviteData, 0, len(ents))

	for _, ent := range ents {
		otherKey, err := s.enc.DecryptUserKey(ent.OtherEncryptedUserKey)
		if err != nil {
			return nil, err
		}

		username, err := s.enc.Decrypt(ent.OtherEncryptedUsername, otherKey)
		if err != nil {
			return nil, err
		}

		firstName, err := s.enc.Decrypt(ent.OtherEncryptedFirstName, otherKey)
		if err != nil {
			return nil, err
		}

		lastName, err := s.enc.Decrypt(ent.OtherEncryptedLastName, otherKey)
		if err != nil {
			return nil, err
		}

		otherImage := ""
		if ent.OtherImage != nil {
			otherImage = *ent.OtherImage
		}

		result = append(result, domain.InviteData{
			ID:             ent.ID.String(),
			OtherUserID:    ent.OtherUserID.String(),
			OtherUsername:  string(username),
			OtherFirstName: string(firstName),
			OtherLastName:  string(lastName),
			OtherImage:     otherImage,
			Status:         ent.Status,
			Note:           ent.Note,
			CreatedAt:      ent.CreatedAt,
			UpdatedAt:      ent.UpdatedAt,
		})
	}

	return result, nil
}

func (s *service) decryptFloat(ciphertext []byte, key []byte) (float64, error) {
	plaintext, err := s.enc.Decrypt(ciphertext, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(plaintext), 64)
}

func (s *service) decryptInt(ciphertext []byte, key []byte) (int, error) {
	plaintext, err := s.enc.Decrypt(ciphertext, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(plaintext))
}
