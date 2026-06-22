package service

import (
	"context"

	accessdb "github.com/findardi/Wadi/server/internal/access/repository/sqlc"
	authdto "github.com/findardi/Wadi/server/internal/auth/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessRepository interface {
	AddMember(ctx context.Context, arg accessdb.AddMemberParams) (accessdb.WorkspaceMember, error)

	DeleteRole(ctx context.Context, id pgtype.UUID) error
	DeleteMember(ctx context.Context, id pgtype.UUID) error

	EditRole(ctx context.Context, arg accessdb.EditRoleParams) (accessdb.WorkspaceRole, error)
	UpdateRole(ctx context.Context, arg accessdb.UpdateRoleParams) (accessdb.WorkspaceMember, error)

	GetRole(ctx context.Context, id pgtype.UUID) (accessdb.WorkspaceRole, error)
	GetRoles(ctx context.Context, workspaceID pgtype.UUID) ([]accessdb.WorkspaceRole, error)
	GetMember(ctx context.Context, id pgtype.UUID) (accessdb.GetMemberRow, error)
	GetMembers(ctx context.Context, workspaceID pgtype.UUID) ([]accessdb.GetMembersRow, error)
	GetMemberByWorkspaceUser(ctx context.Context, arg accessdb.GetMemberByWorkspaceUserParams) (accessdb.WorkspaceMember, error)
	GetWorkspaceInvitation(ctx context.Context, id pgtype.UUID) (accessdb.WorkspaceUserInvitation, error)

	InsertRole(ctx context.Context, arg accessdb.InsertRoleParams) (accessdb.WorkspaceRole, error)
	InsertWorkspaceInvitation(ctx context.Context, arg accessdb.InsertWorkspaceInvitationParams) (accessdb.WorkspaceUserInvitation, error)
	ListWorkspaceInvitations(ctx context.Context, arg accessdb.ListWorkspaceInvitationsParams) ([]accessdb.ListWorkspaceInvitationsRow, error)

	RevokeWorkspaceInvitation(ctx context.Context, id pgtype.UUID) (accessdb.WorkspaceUserInvitation, error)
	ResendInvitation(ctx context.Context, arg accessdb.ResendInvitationParams) (accessdb.WorkspaceUserInvitation, error)
	ReinviteWorkspaceInvitation(ctx context.Context, arg accessdb.ReinviteWorkspaceInvitationParams) (accessdb.WorkspaceUserInvitation, error)

	ExecTx(ctx context.Context, fn func(q *accessdb.Queries) error) error
}

type MailService interface {
	Send(ctx context.Context, to, subject, body string) error
}

type Tokenizer interface {
	Generate() string
	Hash(code string) string
}

type AuthService interface {
	UserExists(ctx context.Context, email string) (authdto.UserResponse, error)
}
