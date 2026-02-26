package config

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/env"
)

type Config struct {
	Port           string
	Url            string
	Schema         string
	MasterKey      []byte
	AllowedOrigins []string
}

func LoadConfig() Config {
	port := env.Get("SERVER_PORT", "8080")
	url := env.Get("DB_URL", "postgresql://user:password@localhost:5432/endurance?sslmode=disable")
	schema := env.Get("DB_SCHEMA", "endurance")
	masterKeyHex := env.Get("MASTER_KEY", "")
	rawOrigins := env.Get("WS_ALLOWED_ORIGINS", "localhost:5173")
	allowedOrigins := strings.Split(rawOrigins, ",")

	masterKey, err := hex.DecodeString(masterKeyHex)
	if err != nil || len(masterKey) != 32 {
		slog.Error("MASTER_KEY must be a 64-character hex string (32 bytes). Generate with: openssl rand -hex 32")
		os.Exit(1)
	}

	return Config{
		Port:           fmt.Sprintf(":%s", port),
		Url:            url,
		Schema:         schema,
		MasterKey:      masterKey,
		AllowedOrigins: allowedOrigins,
	}
}
