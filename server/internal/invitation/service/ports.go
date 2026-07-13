package service

import (
	"context"

	invitationdb "github.com/findardi/Riksa-App/server/internal/invitation/repository/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type InvitationRepo interface {
	AcceptWorkspaceInvitation(ctx context.Context, arg invitationdb.AcceptWorkspaceInvitationParams) (invitationdb.WorkspaceUserInvitation, error)

	GetMyInvitations(ctx context.Context, userID pgtype.UUID) ([]invitationdb.GetMyInvitationsRow, error)
	GetWorkspaceInvitation(ctx context.Context, id pgtype.UUID) (invitationdb.WorkspaceUserInvitation, error)

	RejectWorkspaceInvitation(ctx context.Context, id pgtype.UUID) (invitationdb.WorkspaceUserInvitation, error)

	ExecTx(ctx context.Context, fn func(*invitationdb.Queries) error) error
}
