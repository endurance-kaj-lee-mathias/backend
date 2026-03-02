package config

import (
	"context"
	"log/slog"

	"firebase.google.com/go/v4/messaging"
)

type FirebaseNotifier struct {
	client *messaging.Client
}

func NewFirebaseNotifier(client *messaging.Client) *FirebaseNotifier {
	return &FirebaseNotifier{client: client}
}

func (n *FirebaseNotifier) Notify(ctx context.Context, deviceToken string) error {
	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Daily check-in",
			Body:  "How are you today?",
		},
		Token: deviceToken,
	}

	_, err := n.client.Send(ctx, msg)
	if err != nil {
		slog.Error("sending firebase notification: %w", err)
	}

	return nil
}
