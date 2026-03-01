package config

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/env"
	"google.golang.org/api/option"
)

func NewFirebaseMessagingClient() (*messaging.Client, error) {
	credentialsFile := env.Get("FIREBASE_CREDENTIALS_FILE", "endurance-credentials.json")

	opt := option.WithAuthCredentialsFile(option.ServiceAccount, credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("initializing firebase app: %w", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("initializing firebase messaging client: %w", err)
	}

	return client, nil
}
