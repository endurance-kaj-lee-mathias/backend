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
	Image     string
}

type Invite struct {
	ID        InviteId
	Sender    InviteUser
	Receiver  InviteUser
	Status    InviteStatus
	Note      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewInvite(sender InviteUser, receiver InviteUser, note *string) (Invite, error) {
	if note != nil && len(*note) > 300 {
		return Invite{}, NoteTooLong
	}

	now := time.Now().UTC()
	return Invite{
		ID:        NewInviteId(),
		Sender:    sender,
		Receiver:  receiver,
		Status:    InviteStatusPending,
		Note:      note,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
