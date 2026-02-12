package auth

import (
	"net/http"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	required := make(map[string]struct{}, len(roles))

	for _, r := range roles {
		required[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaims(r.Context())

			if !ok {
				response.WriteError(w, http.StatusUnauthorized, MissingHeader)
				return
			}

			for _, role := range claims.Roles {
				_, ok := required[role]

				if !ok {
					continue
				}

				next.ServeHTTP(w, r)
				return
			}

			response.WriteError(w, http.StatusForbidden, MissingRole)
		})
	}
}
