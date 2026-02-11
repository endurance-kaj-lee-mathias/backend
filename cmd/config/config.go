package config

import (
	"fmt"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/env"
)

type Config struct {
	Port   string
	Url    string
	Schema string
}

func LoadConfig() Config {
	port := env.Get("SERVER_PORT", "8080")
	url := env.Get("DB_URL", "postgresql://user:password@localhost:5432/endurance?sslmode=disable")
	schema := env.Get("DB_SCHEMA", "registrations")

	return Config{
		Port:   fmt.Sprintf(":%s", port),
		Url:    url,
		Schema: schema,
	}
}
