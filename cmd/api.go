package main

import (
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/health"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (server *server) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))

	userHandler := users.Wire(server.db, server.enc)
	supportHandler := support.Wire(server.db, server.enc)
	healthHandler := health.NewHandler(server.db)
	stressHandler := stress.Wire(server.db, server.enc)
	chatsHandler := chats.Wire(server.db, server.enc)
	wsHandler := ws.Wire(server.idp, server.config.AllowedOrigins)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.TokenAuthentication(server.idp))

		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetOrCreate)

			r.Delete("/me", userHandler.DeleteMe)
			// r.Patch("/me/phone-number", userHandler.PatchPhoneNumber)
			r.Patch("/me/introduction", userHandler.PatchIntroduction)
			r.Patch("/me/about", userHandler.PatchAbout)
			r.Patch("/me/image", userHandler.PatchImage)
			r.Put("/me/address", userHandler.UpsertAddress)
			r.Get("/me/address", userHandler.GetAddress)

			r.Get("/search/{username}", userHandler.GetUserByUsername)

			r.Get("/support", supportHandler.GetAll)
			r.Delete("/support/{supportId}", supportHandler.DeleteSupporter)

			r.Get("/{id}", userHandler.GetUser)
			r.Post("/{id}/support", supportHandler.AddMember)
		})

		r.Route("/stress", func(r chi.Router) {
			r.Post("/samples", stressHandler.IngestSample)
		})

		r.Route("/chats", func(r chi.Router) {
			r.Post("/", chatsHandler.StartConversation)
			r.Post("/{conversationId}/messages", chatsHandler.SendMessage)
			r.Get("/{conversationId}/messages", chatsHandler.GetMessages)
		})
	})

	r.Get("/health", healthHandler.Health)
	r.Get("/ws", wsHandler.ServeWS)

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
