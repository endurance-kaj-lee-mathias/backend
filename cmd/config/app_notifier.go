package config

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	wsapp "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
	wsdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/domain"
)

type AppNotifier struct {
	fb *FirebaseNotifier
	ws *wsapp.Manager
}

func NewAppNotifier(fb *FirebaseNotifier, ws *wsapp.Manager) *AppNotifier {
	return &AppNotifier{
		fb: fb,
		ws: ws,
	}
}

func (n *AppNotifier) NotifyInvite(ctx context.Context, userID uuid.UUID, deviceTokens []string) error {
	channel := "notifications:" + userID.String()
	n.ws.Broadcast(channel, wsdomain.OutboundMessage{
		Channel:   channel,
		Content:   "New support request: Someone wants to connect with you as a supporter.",
		CreatedAt: time.Now().UTC(),
	})

	for _, token := range deviceTokens {
		if token != "" {
			_ = n.fb.NotifyInvite(ctx, token)
		}
	}
	return nil
}

func (n *AppNotifier) NotifyInviteAccepted(ctx context.Context, userID uuid.UUID, deviceTokens []string) error {
	channel := "notifications:" + userID.String()
	n.ws.Broadcast(channel, wsdomain.OutboundMessage{
		Channel:   channel,
		Content:   "Support request accepted: Your support request has been accepted.",
		CreatedAt: time.Now().UTC(),
	})

	for _, token := range deviceTokens {
		if token != "" {
			_ = n.fb.NotifyInviteAccepted(ctx, token)
		}
	}
	return nil
}

func (n *AppNotifier) Notify(ctx context.Context, userID uuid.UUID, deviceTokens []string) error {
	channel := "notifications:" + userID.String()
	n.ws.Broadcast(channel, wsdomain.OutboundMessage{
		Channel:   channel,
		Content:   "Daily check-in: How are you today?",
		CreatedAt: time.Now().UTC(),
	})

	for _, token := range deviceTokens {
		if token != "" {
			_ = n.fb.Notify(ctx, token)
		}
	}
	return nil
}

func (n *AppNotifier) NotifyNewMessage(ctx context.Context, userID uuid.UUID, deviceTokens []string) error {
	channel := "notifications:" + userID.String()
	n.ws.Broadcast(channel, wsdomain.OutboundMessage{
		Channel:   channel,
		Content:   "New message: You have received a new message.",
		CreatedAt: time.Now().UTC(),
	})

	for _, token := range deviceTokens {
		if token != "" {
			_ = n.fb.NotifyNewMessage(ctx, token)
		}
	}
	return nil
}

func (n *AppNotifier) NotifyHighStress(ctx context.Context, userID uuid.UUID, deviceTokens []string) error {
	channel := "notifications:" + userID.String()
	n.ws.Broadcast(channel, wsdomain.OutboundMessage{
		Channel:   channel,
		Content:   "High stress detected: Your stress levels are high. Please open the app for support.",
		CreatedAt: time.Now().UTC(),
	})

	for _, token := range deviceTokens {
		if token != "" {
			_ = n.fb.NotifyHighStress(ctx, token)
		}
	}
	return nil
}
