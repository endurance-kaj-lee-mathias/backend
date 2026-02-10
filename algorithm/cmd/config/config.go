package config

import (
	"fmt"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/algorithm-service/internal/env"
)

type Config struct {
	Port string
}

func LoadConfig() Config {
	port := env.Get("SERVER_PORT", "8080")

	return Config{
		Port: fmt.Sprintf(":%s", port),
	}
}
