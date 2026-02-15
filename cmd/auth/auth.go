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

func TokenAuthentication(config config.Idp) func(http.Handler) http.Handler {
	var url = fmt.Sprintf(
		"%s/realms/%s/protocol/openid-connect/certs",
		config.Url,
		config.Realm,
	)

	jwks := initialize(url, config.Refresh)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractToken(r)

			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			var issuer = fmt.Sprintf(
				"%s/realms/%s",
				config.Url,
				config.Realm,
			)

			claims, err := validateToken(token, jwks, issuer, config.Client)

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

func validateToken(tokenString string, jwks *keyfunc.JWKS, issuer, audience string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, jwks.Keyfunc,
		jwt.WithIssuer(issuer),
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

	azp, ok := claims["azp"].(string)
	if !ok || azp != audience {
		return nil, ClaimsInvalid
	}

	return claims, nil
}

func getClaims(context context.Context) (jwt.MapClaims, bool) {
	claims, ok := context.Value(claimsKey).(jwt.MapClaims)
	return claims, ok
}

func GetUserClaims(ctx context.Context) (*Claims, bool) {
	raw, ok := getClaims(ctx)
	if !ok {
		return nil, false
	}

	c := &Claims{}

	if sub, ok := raw["sub"].(string); ok {
		c.Sub = sub
	}

	if email, ok := raw["email"].(string); ok {
		c.Email = email
	}

	if firstName, ok := raw["given_name"].(string); ok {
		c.FirstName = firstName
	}

	if lastName, ok := raw["family_name"].(string); ok {
		c.LastName = lastName
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

	if c.Sub == "" {
		return nil, false
	}

	return c, true
}
