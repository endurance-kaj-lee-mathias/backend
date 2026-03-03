package config

import (
	"strings"
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/env"
)

type Idp struct {
	Url           string
	Issuers       []string
	Realm         string
	Client        string
	Refresh       time.Duration
	AdminUser     string
	AdminPassword string
}

func LoadIdp() Idp {
	url := env.Get("IDP_URL", "http://localhost:8180")
	realm := env.Get("IDP_REALM", "endurance")
	client := env.Get("IDP_CLIENT", "backend")
	adminUser := env.Get("KEYCLOAK_ADMIN_USER", "admin")
	adminPassword := env.Get("KEYCLOAK_ADMIN_PASSWORD", "admin")

	defaultIssuers := "http://localhost:8180,https://10.0.2.2:8443"
	issuersRaw := env.Get("IDP_ISSUERS", defaultIssuers)

	estimatedIssuers := strings.Count(issuersRaw, ",") + 1
	issuers := make([]string, 0, estimatedIssuers)
	for _, s := range strings.Split(issuersRaw, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			issuers = append(issuers, s)
		}
	}
	return Idp{
		Url:           url,
		Issuers:       issuers,
		Realm:         realm,
		Client:        client,
		Refresh:       time.Hour,
		AdminUser:     adminUser,
		AdminPassword: adminPassword,
	}
}
