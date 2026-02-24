package main

import (
	"log/slog"
	"os"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"

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

	db := loadDatabase(cfg.Url, cfg.Schema)
	api := server{cfg, idp, db, enc}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has crashed", "error", err)
		os.Exit(1)
	}
}
