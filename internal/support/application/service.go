package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
	userdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

func convertRoles(roleStrings []string) []userdomain.Role {
	roles := make([]userdomain.Role, len(roleStrings))
	for i, roleStr := range roleStrings {
		roles[i] = userdomain.Role(roleStr)
	}
	return roles
}

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
		roleStrings, err := s.userRoleRead.GetRoles(ctx, member.ID.UUID)
		if err != nil {
			return nil, err
		}

		members[i].Roles = convertRoles(roleStrings)
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
		roleStrings, err := s.userRoleRead.GetRoles(ctx, member.ID.UUID)
		if err != nil {
			return nil, err
		}

		members[i].Roles = convertRoles(roleStrings)
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

func (s *service) DeleteFriend(ctx context.Context, callerID domain.MemberId, friendID domain.MemberId) error {
	if err := s.repo.Delete(ctx, callerID.UUID, friendID.UUID); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, friendID.UUID, callerID.UUID); err != nil {
		return err
	}

	if err := s.authz.RevokeAll(ctx, callerID.UUID, friendID.UUID); err != nil {
		return err
	}

	if err := s.authz.RevokeAll(ctx, friendID.UUID, callerID.UUID); err != nil {
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

	senderRoles, err := s.userRoleRead.GetRoles(ctx, senderID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	receiverRoles, err := s.userRoleRead.GetRoles(ctx, receiverID.UUID)
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

	accepted, err := s.inviteRepo.FindAcceptedBySenderReceiver(ctx, senderID.UUID, receiverID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	if accepted {
		return domain.Invite{}, domain.AlreadyAccepted
	}

	senderUser := domain.InviteUser{ID: senderID, Roles: convertRoles(senderRoles)}
	receiverUser := domain.InviteUser{ID: receiverID, Roles: convertRoles(receiverRoles)}
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

	senderRolesUpdated, err := s.userRoleRead.GetRoles(ctx, senderID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}
	inv.Sender.Roles = convertRoles(senderRolesUpdated)

	receiverRolesUpdated, err := s.userRoleRead.GetRoles(ctx, receiverID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}
	inv.Receiver.Roles = convertRoles(receiverRolesUpdated)

	return inv, nil
}

func (s *service) AcceptInvite(ctx context.Context, callerID domain.MemberId, inviteID domain.InviteId) (domain.Invite, error) {
	ent, err := s.inviteRepo.FindInviteByID(ctx, inviteID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	if ent.ReceiverID != callerID.UUID {
		return domain.Invite{}, domain.NotReceiver
	}

	if err := s.inviteRepo.UpdateInviteStatus(ctx, inviteID.UUID, domain.InviteStatusAccepted); err != nil {
		return domain.Invite{}, err
	}

	if _, err := s.repo.Create(ctx, ent.ReceiverID, ent.SenderID); err != nil {
		return domain.Invite{}, err
	}

	updated, err := s.inviteRepo.FindInviteByID(ctx, inviteID.UUID)
	if err != nil {
		return domain.Invite{}, err
	}

	inv, err := entities.FromInviteEntity(updated, s.enc)
	if err != nil {
		return domain.Invite{}, err
	}

	senderRoles, err := s.userRoleRead.GetRoles(ctx, ent.SenderID)
	if err != nil {
		return domain.Invite{}, err
	}
	inv.Sender.Roles = convertRoles(senderRoles)

	receiverRoles, err := s.userRoleRead.GetRoles(ctx, ent.ReceiverID)
	if err != nil {
		return domain.Invite{}, err
	}
	inv.Receiver.Roles = convertRoles(receiverRoles)

	return inv, nil
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

		senderRoles, err := s.userRoleRead.GetRoles(ctx, ent.SenderID)
		if err != nil {
			return nil, nil, err
		}
		inv.Sender.Roles = convertRoles(senderRoles)

		receiverRoles, err := s.userRoleRead.GetRoles(ctx, ent.ReceiverID)
		if err != nil {
			return nil, nil, err
		}
		inv.Receiver.Roles = convertRoles(receiverRoles)

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
