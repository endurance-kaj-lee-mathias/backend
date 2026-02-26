package application

import (
	"context"
	"log/slog"
	"time"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/infrastructure"
)

type Scheduler struct {
	repo     infrastructure.Repository
	notifier PhoneNotifier
	roleHash string
}

func NewScheduler(repo infrastructure.Repository, notifier PhoneNotifier, veteranRoleHash string) *Scheduler {
	return &Scheduler{repo: repo, notifier: notifier, roleHash: veteranRoleHash}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.run(ctx)
		}
	}
}

func (s *Scheduler) run(ctx context.Context) {
	ids, err := s.repo.FindVeteransWithoutEntryInLast24Hours(ctx, s.roleHash)
	if err != nil {
		slog.Warn("mood scheduler: failed to query veterans", "error", err)
		return
	}

	for _, id := range ids {
		if err := notify(ctx, s.notifier, ""); err != nil {
			slog.Warn("mood scheduler: failed to notify user after retries", "user_id", id, "error", err)
		}
	}
}

func notify(ctx context.Context, notifier PhoneNotifier, deviceToken string) error {
	delays := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second}

	var err error
	for attempt := 0; attempt <= len(delays); attempt++ {
		err = notifier.Notify(ctx, deviceToken)
		if err == nil {
			return nil
		}
		if attempt < len(delays) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delays[attempt]):
			}
		}
	}

	return err
}
