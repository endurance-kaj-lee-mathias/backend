package main

import (
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/health"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood"
	moodapp "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (server *server) mount() (http.Handler, *moodapp.Scheduler) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))

	userHandler := users.Wire(server.db, server.enc, server.kc)
	supportHandler := support.Wire(server.db, server.enc)
	healthHandler := health.NewHandler(server.db, server.messagingClient)
	stressHandler := stress.Wire(server.db, server.enc, server.config.AlgoServiceURL)
	chatsHandler := chats.Wire(server.db, server.enc)
	wsHandler := ws.Wire(server.idp, server.config.AllowedOrigins)
	moodHandler, moodScheduler := mood.Wire(server.db, server.enc, server.notifier)
	calendarHandler := calendar.Wire(server.db, server.config.MinUrgentMinutes)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.TokenAuthentication(server.idp))

		r.Route("/users", func(r chi.Router) {
			r.Get("/me", userHandler.GetOrCreate)

			r.Delete("/me", userHandler.DeleteMe)
			r.Patch("/me/phone-number", userHandler.PatchPhoneNumber)
			r.Patch("/me/introduction", userHandler.PatchIntroduction)
			r.Patch("/me/about", userHandler.PatchAbout)
			r.Patch("/me/image", userHandler.PatchImage)
			r.Put("/me/address", userHandler.UpsertAddress)

			r.Put("/device", userHandler.PutDevice)
			r.Delete("/device", userHandler.DeleteDevice)

			r.Get("/search/{username}", userHandler.GetUserByUsername)

			r.Get("/support", supportHandler.GetAll)
			r.Delete("/support/{supportId}", supportHandler.DeleteSupporter)

			r.Get("/{id}", userHandler.GetUser)
		})

		r.Route("/support-invites", func(r chi.Router) {
			r.Post("/", supportHandler.PostInvite)
			r.Patch("/{inviteId}/accept", supportHandler.AcceptInvite)
			r.Patch("/{inviteId}/decline", supportHandler.DeclineInvite)
			r.Get("/", supportHandler.ListInvites)
		})

		r.Route("/stress", func(r chi.Router) {
			r.Post("/samples", stressHandler.IngestSample)
			r.Get("/scores/latest", stressHandler.GetLatestScore)
		})

		r.Route("/chats", func(r chi.Router) {
			r.Post("/", chatsHandler.StartConversation)
			r.Post("/{conversationId}/messages", chatsHandler.SendMessage)
			r.Get("/{conversationId}/messages", chatsHandler.GetMessages)
		})

		r.Route("/mood", func(r chi.Router) {
			r.Post("/entries", moodHandler.UpsertMoodEntry)
		})

		r.Route("/calendar", func(r chi.Router) {
			r.Post("/slots", calendarHandler.CreateSlot)
			r.Get("/slots", calendarHandler.GetSlots)
			r.Delete("/slots/{id}", calendarHandler.DeleteSlot)
			r.Post("/slots/{id}/book", calendarHandler.BookSlot)
			r.Patch("/appointments/{id}/cancel", calendarHandler.CancelAppointment)
		})
	})

	r.Get("/health", healthHandler.Health)
	r.Get("/ws", wsHandler.ServeWS)

	return r, moodScheduler
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
