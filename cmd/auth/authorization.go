package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
)

type contextKeyResource string
type contextKeyTargetID string

const resourceKey contextKeyResource = "authResource"
const targetIDKey contextKeyTargetID = "authTargetID"

type TargetExtractor func(*http.Request) (uuid.UUID, error)
type UsernameResolver func(context.Context, string) (uuid.UUID, error)

func ExtractTargetFromUsername(resolve UsernameResolver) TargetExtractor {
	return func(r *http.Request) (uuid.UUID, error) {
		targetID, ok := GetTargetID(r.Context())
		if ok {
			return targetID, nil
		}

		username := chi.URLParam(r, "username")
		return resolve(r.Context(), username)
	}
}

func WithResource(resource string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), resourceKey, resource)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetResource(ctx context.Context) string {
	val, _ := ctx.Value(resourceKey).(string)
	return val
}

func SetTargetID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, targetIDKey, id)
}

func GetTargetID(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(targetIDKey).(uuid.UUID)
	return val, ok
}

func RequireSupportRelationship(authService application.Service, extract TargetExtractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaims(r.Context())
			if !ok {
				response.WriteError(w, http.StatusUnauthorized, MissingHeader)
				return
			}

			viewerID, err := uuid.FromString(claims.Sub)
			if err != nil {
				response.WriteError(w, http.StatusBadRequest, ClaimsInvalid)
				return
			}

			targetID, err := extract(r)
			if err != nil {
				if errors.Is(err, TargetNotFound) {
					response.WriteError(w, http.StatusNotFound, err)
					return
				}

				response.WriteError(w, http.StatusBadRequest, err)
				return
			}

			ctx := SetTargetID(r.Context(), targetID)
			r = r.WithContext(ctx)

			if viewerID == targetID {
				next.ServeHTTP(w, r)
				return
			}

			hasSupportRel, err := authService.HasSupportRelationship(r.Context(), viewerID, targetID)
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			if !hasSupportRel {
				response.WriteError(w, http.StatusForbidden, MissingRole)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthorization(authService application.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaims(r.Context())
			if !ok {
				response.WriteError(w, http.StatusUnauthorized, MissingHeader)
				return
			}

			viewerID, err := uuid.FromString(claims.Sub)
			if err != nil {
				response.WriteError(w, http.StatusBadRequest, ClaimsInvalid)
				return
			}

			targetID, ok := GetTargetID(r.Context())
			if !ok {
				response.WriteError(w, http.StatusBadRequest, ClaimsInvalid)
				return
			}

			if viewerID == targetID {
				next.ServeHTTP(w, r)
				return
			}

			resource := GetResource(r.Context())

			allowed, err := authService.IsAllowed(r.Context(), targetID, viewerID, resource)
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			if !allowed {
				response.WriteError(w, http.StatusForbidden, MissingRole)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
