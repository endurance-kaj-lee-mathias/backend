package domain

import "time"

type InviteStatus string

const (
	InviteStatusPending  InviteStatus = "PENDING"
	InviteStatusAccepted InviteStatus = "ACCEPTED"
	InviteStatusDeclined InviteStatus = "DECLINED"
)

type InviteUser struct {
	ID        MemberId
	Username  string
	FirstName string
	LastName  string
}

type Invite struct {
	ID        InviteId
	Sender    InviteUser
	Receiver  InviteUser
	Status    InviteStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewInvite(sender InviteUser, receiver InviteUser) Invite {
	now := time.Now().UTC()
	return Invite{
		ID:        NewInviteId(),
		Sender:    sender,
		Receiver:  receiver,
		Status:    InviteStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
