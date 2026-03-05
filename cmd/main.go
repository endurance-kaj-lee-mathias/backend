package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/keycloak"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("env file was not found", "error", err)
	}

	cfg := config.LoadConfig()
	idp := config.LoadIdp()

	enc, err := encryption.NewService(cfg.MasterKey)
	if err != nil {
		slog.Error("failed to initialize encryption service", "error", err)
		os.Exit(1)
	}

	messagingClient, err := config.NewFirebaseMessagingClient()
	if err != nil {
		slog.Error("failed to initialize firebase messaging client", "error", err)
		os.Exit(1)
	}

	db := loadDatabase(cfg.Url, cfg.Schema)
	kc := keycloak.NewClient(idp.Url, idp.Realm, idp.AdminUser, idp.AdminPassword)
	api := server{
		config:          cfg,
		idp:             idp,
		db:              db,
		enc:             enc,
		kc:              kc,
		notifier:        config.NewFirebaseNotifier(messagingClient),
		messagingClient: messagingClient,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	h, scheduler := api.mount()
	go scheduler.Start(ctx)

	if err := api.run(ctx, h); err != nil {
		slog.Error("server has crashed", "error", err)
		os.Exit(1)
	}
}
