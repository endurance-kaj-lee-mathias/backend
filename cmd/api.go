package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization"
	authzdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/health"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood"
	moodapp "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws"
)

func (server *server) mount() (http.Handler, *moodapp.Scheduler) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))
	r.Use(corsMiddleware(server.config.AllowedOrigins))

	userHandler := users.Wire(server.db, server.enc, server.kc, "mobile", server.idp.WebClientID)
	healthHandler := health.NewHandler(server.db, server.messagingClient)
	stressHandler := stress.Wire(server.db, server.enc, server.config.AlgoServiceURL, server.config.AlgoAPIKey)
	chatsHandler := chats.Wire(server.db, server.enc, server.notifier)
	wsHandler := ws.Wire(server.idp, server.config.AllowedOrigins)
	moodHandler, moodScheduler := mood.Wire(server.db, server.enc, server.notifier)
	calendarHandler := calendar.Wire(server.db, server.enc, server.config.MinUrgentMinutes)
	authzHandler, authzService := authorization.Wire(server.db)
	supportHandler := support.Wire(server.db, server.enc, authzService, server.notifier)
	exportHandler := export.Wire(server.db, server.enc)

	r.Group(func(r chi.Router) {
		r.Use(auth.TokenAuthentication(server.idp))

		r.Route("/users", func(r chi.Router) {
			r.Get("/me", userHandler.GetOrCreate)
			r.Post("/me/register", userHandler.Register)
			r.Get("/me/export", exportHandler.ExportUserData)

			r.Delete("/me", userHandler.DeleteMe)
			r.Patch("/me/phone-number", userHandler.PatchPhoneNumber)
			r.Patch("/me/first-name", userHandler.PatchFirstName)
			r.Patch("/me/last-name", userHandler.PatchLastName)
			r.Patch("/me/introduction", userHandler.PatchIntroduction)
			r.Patch("/me/about", userHandler.PatchAbout)
			r.Patch("/me/image", userHandler.PatchImage)
			r.Patch("/me/privacy", userHandler.PatchPrivacy)
			r.Put("/me/address", userHandler.UpsertAddress)

			r.Put("/device", userHandler.PutDevice)
			r.Delete("/device", userHandler.DeleteDevice)

			r.Get("/support", supportHandler.GetAll)
			r.Delete("/support/{supportId}", supportHandler.DeleteSupporter)

			r.Group(func(r chi.Router) {
				r.Use(auth.WithResource(string(authzdomain.ResourceUserProfile)))
				r.Use(auth.RequireSupportRelationship(authzService, extractTargetFromUsername(userHandler)))
				r.Use(auth.RequireAuthorization(authzService))
				r.Get("/search/{username}", userHandler.GetUserByUsername)
			})

			r.Group(func(r chi.Router) {
				r.Use(auth.WithResource(string(authzdomain.ResourceUserProfile)))
				r.Use(auth.RequireSupportRelationship(authzService, extractTargetFromPathID))
				r.Use(auth.RequireAuthorization(authzService))
				r.Get("/{id}", userHandler.GetUser)
			})
		})

		r.Route("/support", func(r chi.Router) {
			r.Post("/", supportHandler.PostInvite)
			r.Patch("/{inviteId}/accept", supportHandler.AcceptInvite)
			r.Patch("/{inviteId}/decline", supportHandler.DeclineInvite)
			r.Get("/", supportHandler.ListInvites)
		})

		r.Route("/sharing", func(r chi.Router) {
			r.Post("/rules", authzHandler.CreateRule)
			r.Delete("/rules/{id}", authzHandler.DeleteRule)
			r.Get("/rules", authzHandler.ListRules)
			r.Get("/rules/{id}", authzHandler.GetRulesByViewer)
		})

		r.Route("/stress", func(r chi.Router) {
			r.Post("/samples", stressHandler.IngestSample)
			r.Get("/samples/latest", stressHandler.GetLatestSampleTimestamp)
			r.Delete("/samples/me", stressHandler.DeleteMySamples)
			r.Get("/scores/latest", stressHandler.GetLatestScore)

			r.Group(func(r chi.Router) {
				r.Use(auth.WithResource(string(authzdomain.ResourceStressScores)))
				r.Use(auth.RequireSupportRelationship(authzService, extractTargetFromPathID))
				r.Use(auth.RequireAuthorization(authzService))
				r.Get("/scores/{id}/latest", stressHandler.GetLatestScoreByUserID)
			})
		})

		r.Route("/chats", func(r chi.Router) {
			r.Get("/", chatsHandler.GetAllChats)
			r.Post("/", chatsHandler.StartConversation)
			r.Post("/{conversationId}/messages", chatsHandler.SendMessage)
			r.Get("/{conversationId}/messages", chatsHandler.GetMessages)
		})

		r.Route("/mood", func(r chi.Router) {
			r.Post("/entries", moodHandler.UpsertMoodEntry)
			r.Get("/entries/me", moodHandler.GetMyEntries)
			r.Get("/entries/me/today", moodHandler.GetTodayEntry)
			r.Delete("/entries/me/all", moodHandler.DeleteMyEntries)
			r.Put("/entries/{entryId}", moodHandler.UpdateMoodEntry)
			r.Delete("/entries/{entryId}", moodHandler.DeleteMoodEntry)

			r.Group(func(r chi.Router) {
				r.Use(auth.WithResource(string(authzdomain.ResourceMoodEntries)))
				r.Use(auth.RequireSupportRelationship(authzService, extractTargetFromPathID))
				r.Use(auth.RequireAuthorization(authzService))
				r.Get("/entries/{id}", moodHandler.GetEntriesByUserID)
			})
		})

		r.Route("/calendar", func(r chi.Router) {
			r.Get("/me/export", calendarHandler.ExportCalendar)
			r.Get("/me/feed", calendarHandler.FeedCalendar)

			r.Post("/slots", calendarHandler.CreateSlot)
			r.Get("/slots", calendarHandler.GetSlots)
			r.Delete("/slots/me", calendarHandler.DeleteMySlots)
			r.Delete("/slots/{id}", calendarHandler.DeleteSlot)
			r.Post("/slots/{id}/book", calendarHandler.BookSlot)
			r.Patch("/appointments/{id}/cancel", calendarHandler.CancelAppointment)

			r.Group(func(r chi.Router) {
				r.Use(auth.WithResource(string(authzdomain.ResourceCalendar)))
				r.Use(auth.RequireSupportRelationship(authzService, extractTargetFromPathID))
				r.Use(auth.RequireAuthorization(authzService))
				r.Get("/slots/{id}", calendarHandler.GetSlotsByUserID)
			})
		})
	})

	r.Get("/health", healthHandler.Health)
	r.Get("/ws", wsHandler.ServeWS)

	return r, moodScheduler
}

func extractTargetFromPathID(r *http.Request) (uuid.UUID, error) {
	return uuid.FromString(chi.URLParam(r, "id"))
}

func extractTargetFromUsername(userHandler *users.Handler) auth.TargetExtractor {
	return func(r *http.Request) (uuid.UUID, error) {
		targetID, ok := auth.GetTargetID(r.Context())
		if ok {
			return targetID, nil
		}

		username := chi.URLParam(r, "username")
		usr, err := userHandler.ResolveUsername(r.Context(), username)
		if err != nil {
			return uuid.UUID{}, err
		}

		return usr, nil
	}
}

func (server *server) run(ctx context.Context, h http.Handler) error {
	srv := &http.Server{
		Addr:         server.config.Port,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server has started", "port", server.config.Port)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}

		err := <-errCh
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	}
}
