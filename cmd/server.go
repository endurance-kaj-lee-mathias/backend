package main

import (
	"database/sql"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/keycloak"

	"firebase.google.com/go/v4/messaging"
)

type server struct {
	config          config.Config
	idp             config.Idp
	db              *sql.DB
	enc             encryption.Service
	kc              keycloak.Client
	notifier        *config.FirebaseNotifier
	messagingClient *messaging.Client
}
