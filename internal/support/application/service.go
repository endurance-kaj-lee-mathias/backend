package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
	userdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

func (s *service) GetAll(ctx context.Context, id domain.VeteranId) ([]domain.Member, error) {
	ents, err := s.repo.ReadAll(ctx, id.UUID)
	if err != nil {
		return nil, err
	}

	members, err := entities.FromEntities(ents, s.enc)
	if err != nil {
		return nil, err
	}

	for i, member := range members {
		roleStr, err := s.userRoleRead.GetRole(ctx, member.ID.UUID)
		if err != nil {
			return nil, err
		}

		members[i].Role = userdomain.Role(roleStr)
	}

	return members, nil
}

func (s *service) GetAllByMember(ctx context.Context, id domain.MemberId) ([]domain.Member, error) {
	ents, err := s.repo.ReadAllByMember(ctx, id.UUID)
	if err != nil {
		return nil, err
	}

	members, err := entities.FromEntities(ents, s.enc)
	if err != nil {
		return nil, err
	}

	for i, member := range members {
		roleStr, err := s.userRoleRead.GetRole(ctx, member.ID.UUID)
		if err != nil {
			return nil, err
		}

		members[i].Role = userdomain.Role(roleStr)
	}

	return members, nil
}

func (s *service) DeleteSupporter(ctx context.Context, veteranID domain.VeteranId, supportID domain.MemberId) error {
	if err := s.repo.Delete(ctx, veteranID.UUID, supportID.UUID); err != nil {
		return err
	}

	if err := s.authz.RevokeAll(ctx, veteranID.UUID, supportID.UUID); err != nil {
		return err
	}

	return nil
}

func (s *service) SendInvite(ctx context.Context, senderID domain.MemberId, username string, note *string) (domain.Invite, error) {
	receiverUUID, err := s.userRoleRead.FindIDByUsername(ctx, username)
	if err != nil {
		return domain.Invite{}, err
	}

	receiverID, err := domain.ParseMemberId(receiverUUID.String())
	if err != nil {
		return domain.Invite{}, err
	}

	if senderID.UUID == receiverID.UUID {
		return domain.Invite{}, domain.SelfInvite
	}

	senderRoles, err := s.userRoleRead.GetRole(ctx, senderID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	receiverRoles, err := s.userRoleRead.GetRole(ctx, receiverID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	veteranID, err := domain.ParseVeteranId(receiverID.UUID.String())
	if err != nil {
		return domain.Invite{}, err
	}

	if err := domain.ValidateSupportRelationship(receiverRoles, senderRoles, veteranID.UUID.String(), senderID.UUID.String()); err != nil {
		return domain.Invite{}, err
	}

	_, pending, err := s.inviteRepo.FindPendingBySenderReceiver(ctx, senderID.UUID, receiverID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	if pending {
		return domain.Invite{}, domain.DuplicatePendingInvite
	}

	exists, err := s.repo.ExistsRelationship(ctx, senderID.UUID, receiverID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	if exists {
		return domain.Invite{}, domain.AlreadyConnected
	}

	senderUser := domain.InviteUser{ID: senderID, Role: userdomain.Role(senderRoles)}
	receiverUser := domain.InviteUser{ID: receiverID, Role: userdomain.Role(receiverRoles)}
	inv, err := domain.NewInvite(senderUser, receiverUser, note)
	if err != nil {
		return domain.Invite{}, err
	}

	if err := s.inviteRepo.CreateInvite(ctx, inv); err != nil {
		return domain.Invite{}, err
	}

	ent, err := s.inviteRepo.FindInviteByID(ctx, inv.ID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	inv, err = entities.FromInviteEntity(ent, s.enc)
	if err != nil {
		return domain.Invite{}, err
	}

	senderRoleUpdated, err := s.userRoleRead.GetRole(ctx, senderID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}
	inv.Sender.Role = userdomain.Role(senderRoleUpdated)

	receiverRoleUpdated, err := s.userRoleRead.GetRole(ctx, receiverID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}
	inv.Receiver.Role = userdomain.Role(receiverRoleUpdated)

	return inv, nil
}

func (s *service) AcceptInvite(ctx context.Context, callerID domain.MemberId, inviteID domain.InviteId) error {
	ent, err := s.inviteRepo.FindInviteByID(ctx, inviteID.UUID)
	if err != nil {
		return err
	}

	if ent.ReceiverID != callerID.UUID {
		return domain.NotReceiver
	}

	if _, err := s.repo.Create(ctx, ent.ReceiverID, ent.SenderID); err != nil {
		return err
	}

	if err := s.inviteRepo.DeleteInvite(ctx, inviteID.UUID); err != nil {
		return err
	}

	return nil
}

func (s *service) DeclineInvite(ctx context.Context, callerID domain.MemberId, inviteID domain.InviteId) error {
	ent, err := s.inviteRepo.FindInviteByID(ctx, inviteID.UUID)
	if err != nil {
		return err
	}

	if ent.ReceiverID != callerID.UUID {
		return domain.NotReceiver
	}

	return s.inviteRepo.DeleteInvite(ctx, inviteID.UUID)
}

func (s *service) ListInvites(ctx context.Context, callerID domain.MemberId) ([]domain.Invite, []domain.Invite, error) {
	ents, err := s.inviteRepo.ListPendingForUser(ctx, callerID.UUID)
	if err != nil {
		return nil, nil, err
	}

	var incoming, outgoing []domain.Invite

	for _, ent := range ents {
		inv, err := entities.FromInviteEntity(ent, s.enc)
		if err != nil {
			return nil, nil, err
		}

		senderRole, err := s.userRoleRead.GetRole(ctx, ent.SenderID)
		if err != nil {
			return nil, nil, err
		}
		inv.Sender.Role = userdomain.Role(senderRole)

		receiverRole, err := s.userRoleRead.GetRole(ctx, ent.ReceiverID)
		if err != nil {
			return nil, nil, err
		}
		inv.Receiver.Role = userdomain.Role(receiverRole)

		if ent.ReceiverID == callerID.UUID {
			incoming = append(incoming, inv)
		} else {
			outgoing = append(outgoing, inv)
		}
	}

	if incoming == nil {
		incoming = []domain.Invite{}
	}
	if outgoing == nil {
		outgoing = []domain.Invite{}
	}

	return incoming, outgoing, nil
}
