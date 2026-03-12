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
		slog.Error("sending firebase notification", "error", err)
		return err
	}

	return nil
}

func (n *FirebaseNotifier) NotifyInvite(ctx context.Context, deviceToken string) error {
	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "New support request",
			Body:  "Someone wants to connect with you as a supporter.",
		},
		Token: deviceToken,
	}

	_, err := n.client.Send(ctx, msg)
	if err != nil {
		slog.Error("sending firebase invite notification", "error", err)
		return err
	}

	return nil
}

func (n *FirebaseNotifier) NotifyInviteAccepted(ctx context.Context, deviceToken string) error {
	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Support request accepted",
			Body:  "Your support request has been accepted.",
		},
		Token: deviceToken,
	}

	_, err := n.client.Send(ctx, msg)
	if err != nil {
		slog.Error("sending firebase invite accepted notification", "error", err)
		return err
	}

	return nil
}

func (n *FirebaseNotifier) NotifyNewMessage(ctx context.Context, deviceToken string) error {
	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "New message",
			Body:  "You have received a new message.",
		},
		Token: deviceToken,
	}

	_, err := n.client.Send(ctx, msg)
	if err != nil {
		slog.Error("sending firebase message notification", "error", err)
		return err
	}

	return nil
}
