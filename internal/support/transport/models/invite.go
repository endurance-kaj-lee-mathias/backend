package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
)

type CreateInviteRequest struct {
	ReceiverID string `json:"receiverId"`
}

type InviteUserModel struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}

type InviteModel struct {
	ID        uuid.UUID       `json:"id"`
	Sender    InviteUserModel `json:"sender"`
	Receiver  InviteUserModel `json:"receiver"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type InviteListModel struct {
	Incoming []InviteModel `json:"incoming"`
	Outgoing []InviteModel `json:"outgoing"`
}

func ToInviteModel(inv domain.Invite) InviteModel {
	return InviteModel{
		ID: inv.ID.UUID,
		Sender: InviteUserModel{
			ID:        inv.Sender.ID.UUID,
			Username:  inv.Sender.Username,
			FirstName: inv.Sender.FirstName,
			LastName:  inv.Sender.LastName,
		},
		Receiver: InviteUserModel{
			ID:        inv.Receiver.ID.UUID,
			Username:  inv.Receiver.Username,
			FirstName: inv.Receiver.FirstName,
			LastName:  inv.Receiver.LastName,
		},
		Status:    string(inv.Status),
		CreatedAt: inv.CreatedAt,
		UpdatedAt: inv.UpdatedAt,
	}
}

func ToInviteModels(invs []domain.Invite) []InviteModel {
	out := make([]InviteModel, 0, len(invs))
	for _, inv := range invs {
		out = append(out, ToInviteModel(inv))
	}
	return out
}
