package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
)

type InviteEntity struct {
	ID                        uuid.UUID `db:"id"`
	SenderID                  uuid.UUID `db:"sender_id"`
	SenderEncryptedUsername   []byte    `db:"sender_encrypted_username"`
	SenderEncryptedFirst      []byte    `db:"sender_encrypted_first_name"`
	SenderEncryptedLast       []byte    `db:"sender_encrypted_last_name"`
	SenderEncryptedUserKey    []byte    `db:"sender_encrypted_user_key"`
	SenderImage               *string   `db:"sender_image"`
	ReceiverID                uuid.UUID `db:"receiver_id"`
	ReceiverEncryptedUsername []byte    `db:"receiver_encrypted_username"`
	ReceiverEncryptedFirst    []byte    `db:"receiver_encrypted_first_name"`
	ReceiverEncryptedLast     []byte    `db:"receiver_encrypted_last_name"`
	ReceiverEncryptedUserKey  []byte    `db:"receiver_encrypted_user_key"`
	ReceiverImage             *string   `db:"receiver_image"`
	Status                    string    `db:"status"`
	Note                      *string   `db:"note"`
	CreatedAt                 time.Time `db:"created_at"`
	UpdatedAt                 time.Time `db:"updated_at"`
}

func decryptInviteUser(id uuid.UUID, encKey, encUsername, encFirst, encLast []byte, image *string, enc encryption.Service) (domain.InviteUser, error) {
	userKey, err := enc.DecryptUserKey(encKey)
	if err != nil {
		return domain.InviteUser{}, err
	}

	usernameBytes, err := enc.Decrypt(encUsername, userKey)
	if err != nil {
		return domain.InviteUser{}, err
	}

	firstBytes, err := enc.Decrypt(encFirst, userKey)
	if err != nil {
		return domain.InviteUser{}, err
	}

	lastBytes, err := enc.Decrypt(encLast, userKey)
	if err != nil {
		return domain.InviteUser{}, err
	}

	memberId, err := domain.ParseMemberId(id.String())
	if err != nil {
		return domain.InviteUser{}, err
	}

	return domain.InviteUser{
		ID:        memberId,
		Username:  string(usernameBytes),
		FirstName: string(firstBytes),
		LastName:  string(lastBytes),
		Image:     derefString(image),
	}, nil
}

func FromInviteEntity(ent InviteEntity, enc encryption.Service) (domain.Invite, error) {
	sender, err := decryptInviteUser(
		ent.SenderID,
		ent.SenderEncryptedUserKey, ent.SenderEncryptedUsername,
		ent.SenderEncryptedFirst, ent.SenderEncryptedLast,
		ent.SenderImage,
		enc,
	)
	if err != nil {
		return domain.Invite{}, err
	}

	receiver, err := decryptInviteUser(
		ent.ReceiverID,
		ent.ReceiverEncryptedUserKey, ent.ReceiverEncryptedUsername,
		ent.ReceiverEncryptedFirst, ent.ReceiverEncryptedLast,
		ent.ReceiverImage,
		enc,
	)
	if err != nil {
		return domain.Invite{}, err
	}

	inviteId, err := domain.ParseInviteId(ent.ID.String())
	if err != nil {
		return domain.Invite{}, err
	}

	return domain.Invite{
		ID:        inviteId,
		Sender:    sender,
		Receiver:  receiver,
		Status:    domain.InviteStatus(ent.Status),
		Note:      ent.Note,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}, nil
}

func FromInviteEntities(ents []InviteEntity, enc encryption.Service) ([]domain.Invite, error) {
	out := make([]domain.Invite, 0, len(ents))
	for _, ent := range ents {
		inv, err := FromInviteEntity(ent, enc)
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, nil
}
