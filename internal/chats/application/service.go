package application

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/infrastructure/entities"
	wsdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/domain"
)

func (s *service) GetOrCreateConversation(ctx context.Context, callerID, participantID uuid.UUID) (domain.Conversation, error) {
	ok, err := s.repo.CheckSupportRelationship(ctx, callerID, participantID)
	if err != nil {
		return domain.Conversation{}, err
	}
	if !ok {
		return domain.Conversation{}, NoSupportRelationship
	}

	ent, err := s.repo.FindConversation(ctx, callerID, participantID)
	if err == nil {
		return domain.NewConversation(
			domain.ConversationId{UUID: ent.ID},
			[]uuid.UUID{callerID, participantID},
			ent.CreatedAt,
		), nil
	}

	if !errors.Is(err, infrastructure.ConversationNotFound) {
		return domain.Conversation{}, err
	}

	convID := uuid.Must(uuid.NewV4())
	convEnt := entities.ConversationEntity{
		ID:        convID,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateConversation(ctx, convEnt); err != nil {
		return domain.Conversation{}, err
	}

	convKey, err := s.enc.GenerateUserEncryptionKey()
	if err != nil {
		return domain.Conversation{}, err
	}

	for _, userID := range []uuid.UUID{callerID, participantID} {
		encUserKey, err := s.repo.GetUserEncryptedKey(ctx, userID)
		if err != nil {
			return domain.Conversation{}, err
		}

		userKey, err := s.enc.DecryptUserKey(encUserKey)
		if err != nil {
			return domain.Conversation{}, err
		}

		encConvKey, err := s.enc.Encrypt(convKey, userKey)
		if err != nil {
			return domain.Conversation{}, err
		}

		if err := s.repo.SaveParticipantKey(ctx, entities.ParticipantKeyEntity{
			ConversationID:           convID,
			UserID:                   userID,
			EncryptedConversationKey: encConvKey,
		}); err != nil {
			return domain.Conversation{}, err
		}
	}

	return domain.NewConversation(
		domain.ConversationId{UUID: convID},
		[]uuid.UUID{callerID, participantID},
		convEnt.CreatedAt,
	), nil
}

func (s *service) SendMessage(ctx context.Context, conversationID uuid.UUID, senderID uuid.UUID, content string) (domain.Message, error) {
	convKey, err := s.decryptConversationKey(ctx, conversationID, senderID)
	if err != nil {
		return domain.Message{}, err
	}

	encContent, err := s.enc.Encrypt([]byte(content), convKey)
	if err != nil {
		return domain.Message{}, err
	}

	msgID := uuid.Must(uuid.NewV4())
	now := time.Now().UTC()

	ent := entities.MessageEntity{
		ID:               msgID,
		ConversationID:   conversationID,
		SenderID:         senderID,
		EncryptedContent: encContent,
		CreatedAt:        now,
	}

	if err := s.repo.CreateMessage(ctx, ent); err != nil {
		return domain.Message{}, err
	}

	channel := conversationChannel(conversationID)

	username := ""
	if encUsername, err := s.repo.GetEncryptedUsername(context.Background(), senderID); err != nil {
		slog.Warn("failed to fetch encrypted username for websocket broadcast", "user_id", senderID, "error", err)
	} else {
		encUserKey, err := s.repo.GetUserEncryptedKey(context.Background(), senderID)
		if err != nil {
			slog.Warn("failed to fetch encrypted user key for websocket broadcast", "user_id", senderID, "error", err)
		} else {
			userKey, err := s.enc.DecryptUserKey(encUserKey)
			if err != nil {
				slog.Warn("failed to decrypt user key for websocket broadcast", "user_id", senderID, "error", err)
			} else {
				unameBytes, err := s.enc.Decrypt(encUsername, userKey)
				if err != nil {
					slog.Warn("failed to decrypt username for websocket broadcast", "user_id", senderID, "error", err)
				} else {
					username = string(unameBytes)
				}
			}
		}
	}

	s.broadcaster.Broadcast(channel, wsdomain.OutboundMessage{
		Channel:   channel,
		SenderID:  senderID.String(),
		Username:  username,
		Content:   content,
		CreatedAt: now,
	})

	go s.notifyMessage(conversationID, senderID)

	return domain.NewMessage(
		domain.MessageId{UUID: msgID},
		domain.ConversationId{UUID: conversationID},
		senderID,
		username,
		content,
		now,
	), nil
}

func (s *service) ListConversationIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.ListConversationIDsByUserID(ctx, userID)
}

func (s *service) notifyMessage(conversationID, senderID uuid.UUID) {
	ctx := context.Background()

	recipients, err := s.repo.FindOtherParticipants(ctx, conversationID, senderID)
	if err != nil {
		slog.Warn("failed to fetch other participants for message notification", "conversation_id", conversationID, "error", err)
		return
	}

	for _, recipientID := range recipients {
		tokens, err := s.repo.FindDeviceTokensByUserID(ctx, recipientID)
		if err != nil {
			slog.Warn("failed to fetch device tokens for message notification", "user_id", recipientID, "error", err)
			continue
		}

		for _, token := range tokens {
			if token == "" {
				continue
			}

			if err := s.notifier.NotifyNewMessage(ctx, token); err != nil {
				slog.Warn("failed to send message notification", "user_id", recipientID, "error", err)
			}
		}
	}
}

func (s *service) GetMessages(ctx context.Context, conversationID uuid.UUID, callerID uuid.UUID, limit, offset int) ([]domain.Message, error) {
	convKey, err := s.decryptConversationKey(ctx, conversationID, callerID)
	if err != nil {
		return nil, err
	}

	ents, err := s.repo.GetMessages(ctx, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}

	return entities.FromMessageEntities(ents, convKey, s.enc)
}

func (s *service) GetAllChats(ctx context.Context, userID uuid.UUID) ([]domain.ConversationSummary, error) {
	summaryEnts, err := s.repo.GetConversationSummaries(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(summaryEnts) == 0 {
		return []domain.ConversationSummary{}, nil
	}

	callerUserKey, err := s.enc.DecryptUserKey(summaryEnts[0].CallerEncryptedUserKey)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.ConversationSummary, 0, len(summaryEnts))

	for _, ent := range summaryEnts {
		convKey, err := s.enc.Decrypt(ent.CallerEncryptedConversationKey, callerUserKey)
		if err != nil {
			return nil, err
		}

		var latestMessage *string

		if len(ent.LatestEncryptedContent) > 0 {
			contentBytes, err := s.enc.Decrypt(ent.LatestEncryptedContent, convKey)
			if err != nil {
				return nil, err
			}
			content := string(contentBytes)
			latestMessage = &content
		}

		otherUserKey, err := s.enc.DecryptUserKey(ent.OtherEncryptedUserKey)
		if err != nil {
			return nil, err
		}

		usernameBytes, err := s.enc.Decrypt(ent.OtherEncryptedUsername, otherUserKey)
		if err != nil {
			return nil, err
		}

		firstNameBytes, err := s.enc.Decrypt(ent.OtherEncryptedFirstName, otherUserKey)
		if err != nil {
			return nil, err
		}

		lastNameBytes, err := s.enc.Decrypt(ent.OtherEncryptedLastName, otherUserKey)
		if err != nil {
			return nil, err
		}

		image := ""
		if ent.OtherImage != nil {
			image = *ent.OtherImage
		}

		summaries = append(summaries, domain.ConversationSummary{
			ConversationID:        domain.ConversationId{UUID: ent.ConversationID},
			OtherUserID:           ent.OtherUserID,
			Username:              string(usernameBytes),
			FirstName:             string(firstNameBytes),
			LastName:              string(lastNameBytes),
			Image:                 image,
			LatestMessage:         latestMessage,
			LatestMessageSenderID: ent.LatestSenderID,
			LatestMessageAt:       ent.LatestMessageAt,
		})
	}

	return summaries, nil
}

func (s *service) decryptConversationKey(ctx context.Context, conversationID, userID uuid.UUID) ([]byte, error) {
	pkEnt, err := s.repo.GetParticipantKey(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}

	userKey, err := s.enc.DecryptUserKey(pkEnt.EncryptedUserKey)
	if err != nil {
		return nil, err
	}

	convKey, err := s.enc.Decrypt(pkEnt.EncryptedConversationKey, userKey)
	if err != nil {
		return nil, err
	}

	return convKey, nil
}

func conversationChannel(conversationID uuid.UUID) string {
	return "conversation:" + conversationID.String()
}
