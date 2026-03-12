package application

import (
	"context"
	"sync"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
	userdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"golang.org/x/sync/errgroup"
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

	g, gCtx := errgroup.WithContext(ctx)
	roles := make([]string, len(members))

	for i, member := range members {
		i, member := i, member
		g.Go(func() error {
			roleStr, err := s.userRoleRead.GetRole(gCtx, member.ID.UUID)
			if err != nil {
				return err
			}
			roles[i] = roleStr
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for i := range members {
		members[i].Role = userdomain.Role(roles[i])
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

	g, gCtx := errgroup.WithContext(ctx)
	roles := make([]string, len(members))

	for i, member := range members {
		i, member := i, member
		g.Go(func() error {
			roleStr, err := s.userRoleRead.GetRole(gCtx, member.ID.UUID)
			if err != nil {
				return err
			}
			roles[i] = roleStr
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for i := range members {
		members[i].Role = userdomain.Role(roles[i])
	}

	return members, nil
}

func (s *service) DeleteSupporter(ctx context.Context, callerID domain.MemberId, otherID domain.MemberId) error {
	if err := s.repo.Delete(ctx, callerID.UUID, otherID.UUID); err != nil {
		return err
	}

	if err := s.authz.RevokeAll(ctx, callerID.UUID, otherID.UUID); err != nil {
		return err
	}

	if err := s.authz.RevokeAll(ctx, otherID.UUID, callerID.UUID); err != nil {
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

	var senderRoles, receiverRoles string

	{
		g, gCtx := errgroup.WithContext(ctx)

		g.Go(func() error {
			var err error
			senderRoles, err = s.userRoleRead.GetRole(gCtx, senderID.UUID)
			return err
		})

		g.Go(func() error {
			var err error
			receiverRoles, err = s.userRoleRead.GetRole(gCtx, receiverID.UUID)
			return err
		})

		if err := g.Wait(); err != nil {
			return domain.Invite{}, err
		}
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

	{
		g, gCtx := errgroup.WithContext(ctx)
		var senderRoleUpdated, receiverRoleUpdated string

		g.Go(func() error {
			var err error
			senderRoleUpdated, err = s.userRoleRead.GetRole(gCtx, senderID.UUID)
			return err
		})

		g.Go(func() error {
			var err error
			receiverRoleUpdated, err = s.userRoleRead.GetRole(gCtx, receiverID.UUID)
			return err
		})

		if err := g.Wait(); err != nil {
			return domain.Invite{}, err
		}

		inv.Sender.Role = userdomain.Role(senderRoleUpdated)
		inv.Receiver.Role = userdomain.Role(receiverRoleUpdated)
	}

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

	type result struct {
		inv        domain.Invite
		isIncoming bool
	}

	results := make([]result, len(ents))
	g, gCtx := errgroup.WithContext(ctx)

	for i, ent := range ents {
		i, ent := i, ent
		g.Go(func() error {
			inv, err := entities.FromInviteEntity(ent, s.enc)
			if err != nil {
				return err
			}

			var mu sync.Mutex
			rg, rgCtx := errgroup.WithContext(gCtx)
			var senderRole, receiverRole string

			rg.Go(func() error {
				r, err := s.userRoleRead.GetRole(rgCtx, ent.SenderID)
				if err != nil {
					return err
				}
				mu.Lock()
				senderRole = r
				mu.Unlock()
				return nil
			})

			rg.Go(func() error {
				r, err := s.userRoleRead.GetRole(rgCtx, ent.ReceiverID)
				if err != nil {
					return err
				}
				mu.Lock()
				receiverRole = r
				mu.Unlock()
				return nil
			})

			if err := rg.Wait(); err != nil {
				return err
			}

			inv.Sender.Role = userdomain.Role(senderRole)
			inv.Receiver.Role = userdomain.Role(receiverRole)
			results[i] = result{inv: inv, isIncoming: ent.ReceiverID == callerID.UUID}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	var incoming, outgoing []domain.Invite

	for _, r := range results {
		if r.isIncoming {
			incoming = append(incoming, r.inv)
		} else {
			outgoing = append(outgoing, r.inv)
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
