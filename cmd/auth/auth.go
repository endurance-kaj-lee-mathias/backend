package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/config"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

type contextKey string

const claimsKey contextKey = "claims"

func Authenticate(config config.Idp) func(*http.Request) (jwt.MapClaims, error) {
	url := fmt.Sprintf(
		"%s/realms/%s/protocol/openid-connect/certs",
		config.Url,
		config.Realm,
	)

	jwks := initialize(url, config.Refresh)

	issuers := make([]string, len(config.Issuers))
	for i, iss := range config.Issuers {
		issuers[i] = fmt.Sprintf("%s/realms/%s", iss, config.Realm)
	}

	return func(r *http.Request) (jwt.MapClaims, error) {
		token, err := extractToken(r)
		if err != nil {
			return nil, err
		}

		return validateToken(token, jwks, issuers, config.Client)
	}
}

func AuthenticateWSClaims(config config.Idp) func(*http.Request) (*Claims, error) {
	raw := Authenticate(config)
	return func(r *http.Request) (*Claims, error) {
		mapClaims, err := raw(r)
		if err != nil {
			return nil, err
		}
		c := claimsFromMap(mapClaims)
		if c.Sub == "" {
			return nil, ClaimsInvalid
		}
		return c, nil
	}
}

func AuthenticateClaims(config config.Idp) func(*http.Request) (*Claims, error) {
	raw := Authenticate(config)
	return func(r *http.Request) (*Claims, error) {
		mapClaims, err := raw(r)
		if err != nil {
			return nil, err
		}
		c := claimsFromMap(mapClaims)
		if c.Sub == "" {
			return nil, ClaimsInvalid
		}
		return c, nil
	}
}

func TokenAuthentication(config config.Idp) func(http.Handler) http.Handler {
	authenticate := Authenticate(config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := authenticate(r)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func initialize(url string, refresh time.Duration) *keyfunc.JWKS {
	jwks, err := keyfunc.Get(url, keyfunc.Options{
		RefreshInterval: refresh,
		RefreshErrorHandler: func(error error) {
			log.Println(error.Error())
		},
	})

	if err != nil {
		panic(err.Error())
	}

	return jwks
}

func extractToken(request *http.Request) (string, error) {
	const prefix = "Bearer "
	header := request.Header.Get("Authorization")

	if header == "" {
		return "", MissingHeader
	}

	if !strings.HasPrefix(header, prefix) {
		return "", HeaderInvalid
	}

	token := strings.TrimPrefix(header, prefix)

	if token == "" {
		return "", HeaderInvalid
	}

	return token, nil
}

func validateToken(tokenString string, jwks *keyfunc.JWKS, issuers []string, audience string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, jwks.Keyfunc,
		jwt.WithAudience(audience),
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, TokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ClaimsInvalid
	}

	iss, err := claims.GetIssuer()
	if err != nil {
		return nil, err
	}

	validIssuer := false
	for _, allowed := range issuers {
		if iss == allowed {
			validIssuer = true
			break
		}
	}
	if !validIssuer {
		return nil, IssuerInvalid
	}

	return claims, nil
}

func getClaims(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.MapClaims)
	return claims, ok
}

func claimsFromMap(raw jwt.MapClaims) *Claims {
	c := &Claims{}
	if sub, ok := raw["sub"].(string); ok {
		c.Sub = sub
	}
	if email, ok := raw["email"].(string); ok {
		c.Email = email
	}
	if username, ok := raw["preferred_username"].(string); ok {
		c.Username = username
	}
	if firstName, ok := raw["given_name"].(string); ok {
		c.FirstName = firstName
	}
	if lastName, ok := raw["family_name"].(string); ok {
		c.LastName = lastName
	}
	if phone, ok := raw["phoneNumber"].(string); ok {
		c.PhoneNumber = phone
	}
	if addr, ok := raw["address"].(map[string]any); ok {
		if v, ok := addr["street_address"].(string); ok {
			c.Address.StreetAddress = v
		}
		if v, ok := addr["locality"].(string); ok {
			c.Address.Locality = v
		}
		if v, ok := addr["region"].(string); ok {
			c.Address.Region = v
		}
		if v, ok := addr["postal_code"].(string); ok {
			c.Address.PostalCode = v
		}
		if v, ok := addr["country"].(string); ok {
			c.Address.Country = v
		}
	}
	if ra, ok := raw["realm_access"].(map[string]any); ok {
		if roles, ok := ra["roles"].([]any); ok {
			for _, r := range roles {
				if role, ok := r.(string); ok {
					c.Roles = append(c.Roles, role)
				}
			}
		}
	}
	return c
}

func GetUserClaims(ctx context.Context) (*Claims, bool) {
	raw, ok := getClaims(ctx)
	if !ok {
		return nil, false
	}

	c := claimsFromMap(raw)
	if c.Sub == "" {
		return nil, false
	}

	return c, true
}
