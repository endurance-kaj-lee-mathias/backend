package main

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("env file was not found", "error", err)
	}

	cfg := config.LoadConfig()
	db := loadDatabase(cfg.Url, cfg.Schema)
	api := server{cfg, db}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has crashed", "error", err)
		os.Exit(1)
	}
}
