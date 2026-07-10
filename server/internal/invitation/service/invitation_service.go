package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/findardi/Riksa-App/server/internal/invitation/dto"
	invitationdb "github.com/findardi/Riksa-App/server/internal/invitation/repository/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvitationNotFound   = errors.New("invitation not found")
	ErrInvitationForbidden  = errors.New("invitation does not belong to this user")
	ErrInvitationNotPending = errors.New("invitation is no longer pending")
	ErrInvitationExpired    = errors.New("invitation has expired")
)

type InvitationService struct {
	repo InvitationRepo
}

func NewInvitationService(repo InvitationRepo) *InvitationService {
	return &InvitationService{
		repo: repo,
	}
}

func uuidString(u pgtype.UUID) string {
	v, err := u.Value()
	if err != nil || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *InvitationService) GetListInvitations(ctx context.Context, userID string) ([]dto.GetMyInvitationsRow, error) {
	invitations := []dto.GetMyInvitationsRow{}
	var uID pgtype.UUID
	if err := uID.Scan(userID); err != nil {
		return invitations, fmt.Errorf("parse user id: %w", err)
	}

	invts, err := s.repo.GetMyInvitations(ctx, uID)
	if err != nil {
		return invitations, fmt.Errorf("get invitations: %w", err)
	}

	for _, inv := range invts {
		invitations = append(invitations, dto.GetMyInvitationsRow{
			ID:            uuidString(inv.ID),
			WorkspaceName: deref(inv.WorkspaceName),
			RoleName:      deref(inv.RoleName),
			InvitedBy:     deref(inv.InvitedName),
			ExpiresAt:     inv.ExpiresAt.Time,
			Status:        inv.Status,
		})
	}

	return invitations, nil
}

func (s *InvitationService) AcceptInvitation(ctx context.Context, invitationID, userID string) error {
	var invID, uID pgtype.UUID
	if err := invID.Scan(invitationID); err != nil {
		return fmt.Errorf("parse invitation id: %w", err)
	}
	if err := uID.Scan(userID); err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}

	return s.repo.ExecTx(ctx, func(q *invitationdb.Queries) error {
		inv, err := q.GetWorkspaceInvitation(ctx, invID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvitationNotFound
		}
		if err != nil {
			return fmt.Errorf("get invitation: %w", err)
		}

		if uuidString(inv.UserID) != userID {
			return ErrInvitationForbidden
		}
		if inv.Status != "pending" {
			return ErrInvitationNotPending
		}
		if inv.ExpiresAt.Valid && inv.ExpiresAt.Time.Before(time.Now()) {
			return ErrInvitationExpired
		}

		if _, err := q.AcceptWorkspaceInvitation(ctx, invitationdb.AcceptWorkspaceInvitationParams{
			ID:     invID,
			UserID: uID,
		}); err != nil {
			return fmt.Errorf("accept invitation: %w", err)
		}

		if err := q.InsertMember(ctx, invitationdb.InsertMemberParams{
			WorkspaceID: inv.WorkspaceID,
			UserID:      uID,
			RoleID:      inv.RoleID,
		}); err != nil {
			return fmt.Errorf("add member: %w", err)
		}

		return nil
	})
}

func (s *InvitationService) RejectInvitation(ctx context.Context, invitationID, userID string) error {
	var invID pgtype.UUID
	if err := invID.Scan(invitationID); err != nil {
		return fmt.Errorf("parse invitation id: %w", err)
	}

	inv, err := s.repo.GetWorkspaceInvitation(ctx, invID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvitationNotFound
	}
	if err != nil {
		return fmt.Errorf("get invitation: %w", err)
	}

	if uuidString(inv.UserID) != userID {
		return ErrInvitationForbidden
	}

	if _, err := s.repo.RejectWorkspaceInvitation(ctx, invID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvitationNotPending
		}
		return fmt.Errorf("reject invitation: %w", err)
	}

	return nil
}
