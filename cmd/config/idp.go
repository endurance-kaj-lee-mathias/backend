package config

import (
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/env"
)

type Idp struct {
	Url     string
	Realm   string
	Client  string
	Refresh time.Duration
}

func LoadIdp() Idp {
	url := env.Get("IDP_URL", "http://localhost:8180")
	realm := env.Get("IDP_REALM", "endurance")
	client := env.Get("IDP_CLIENT", "backend")

	return Idp{
		Url:     url,
		Realm:   realm,
		Client:  client,
		Refresh: time.Hour,
	}
}
