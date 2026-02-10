package config

import (
	"endurance/internal/env"
	"fmt"
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
