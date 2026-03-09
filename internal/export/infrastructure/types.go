package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

type Repository interface {
	GetUserWithAddress(ctx context.Context, userID uuid.UUID) (entities.UserExportEntity, *entities.AddressExportEntity, error)
	GetStressSamples(ctx context.Context, userID uuid.UUID) ([]entities.StressSampleExportEntity, error)
	GetStressScores(ctx context.Context, userID uuid.UUID) ([]entities.StressScoreExportEntity, error)
	GetMoodEntries(ctx context.Context, userID uuid.UUID) ([]entities.MoodEntryExportEntity, error)
	GetConversationsAndMessages(ctx context.Context, userID uuid.UUID) ([]entities.MessageExportEntity, error)
	GetAppointmentsWithSlots(ctx context.Context, userID uuid.UUID) ([]entities.AppointmentExportEntity, error)
	GetSlotsAsProvider(ctx context.Context, providerID uuid.UUID) ([]entities.SlotExportEntity, error)
	GetSupporters(ctx context.Context, veteranID uuid.UUID) ([]entities.SupportMemberExportEntity, error)
	GetSupportedVeterans(ctx context.Context, supportID uuid.UUID) ([]entities.SupportMemberExportEntity, error)
	GetAuthorizationRulesAsOwner(ctx context.Context, userID uuid.UUID) ([]entities.AuthorizationRuleExportEntity, error)
	GetAuthorizationRulesAsViewer(ctx context.Context, userID uuid.UUID) ([]entities.AuthorizationRuleExportEntity, error)
	GetDevices(ctx context.Context, userID uuid.UUID) ([]entities.DeviceExportEntity, error)
	GetSentInvites(ctx context.Context, userID uuid.UUID) ([]entities.InviteExportEntity, error)
	GetReceivedInvites(ctx context.Context, userID uuid.UUID) ([]entities.InviteExportEntity, error)
}

type repository struct {
	db  *sql.DB
	enc encryption.Service
}

func NewRepository(db *sql.DB, enc encryption.Service) Repository {
	return &repository{db: db, enc: enc}
}
