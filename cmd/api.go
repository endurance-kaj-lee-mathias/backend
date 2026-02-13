package main

import (
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/health"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (server *server) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))

	handler := message.Wire(server.db)
	userHandler := users.Wire(server.db)
	healthHandler := health.NewHandler(server.db)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.TokenAuthentication(server.idp))
		r.Get("/hello", handler.GetMessage)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireRoles("admin"))
			r.Get("/hello-admin", handler.GetMessage)
		})

		r.Get("/hello-token", handler.GetMessage)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Post("/{veteranId}/support", userHandler.AddSupportMember)
			r.Get("/{veteranId}/support", userHandler.ListSupportMembers)
		})
	})

	r.Get("/hello-public", handler.GetMessage)
	r.Get("/health", healthHandler.Health)
	return r
}

func (server *server) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         server.config.Port,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server has started", "port", server.config.Port)
	return srv.ListenAndServe()
}
