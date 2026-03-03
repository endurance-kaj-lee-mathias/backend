package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/stress/domain"
)

type StressScoreEntity struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	Score        float64   `db:"score"`
	Category     string    `db:"category"`
	ModelVersion string    `db:"model_version"`
	ComputedAt   time.Time `db:"computed_at"`
}

func ScoreToEntity(score domain.StressScore) StressScoreEntity {
	return StressScoreEntity{
		ID:           score.ID.UUID,
		UserID:       score.UserID,
		Score:        score.Score,
		Category:     score.Category,
		ModelVersion: score.ModelVersion,
		ComputedAt:   score.ComputedAt,
	}
}

func ScoreFromEntity(ent StressScoreEntity) domain.StressScore {
	return domain.StressScore{
		ID:           domain.ScoreId{UUID: ent.ID},
		UserID:       ent.UserID,
		Score:        ent.Score,
		Category:     ent.Category,
		ModelVersion: ent.ModelVersion,
		ComputedAt:   ent.ComputedAt,
	}
}
