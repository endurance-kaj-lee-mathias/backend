package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
)

type CreateInviteRequest struct {
	Username string  `json:"username"`
	Note     *string `json:"note,omitempty"`
}

type InviteUserModel struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Image     string    `json:"image"`
}

type InviteModel struct {
	ID        uuid.UUID       `json:"id"`
	Sender    InviteUserModel `json:"sender"`
	Receiver  InviteUserModel `json:"receiver"`
	Status    string          `json:"status"`
	Note      *string         `json:"note,omitempty"`
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
			Image:     inv.Sender.Image,
		},
		Receiver: InviteUserModel{
			ID:        inv.Receiver.ID.UUID,
			Username:  inv.Receiver.Username,
			FirstName: inv.Receiver.FirstName,
			LastName:  inv.Receiver.LastName,
			Image:     inv.Receiver.Image,
		},
		Status:    string(inv.Status),
		Note:      inv.Note,
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
