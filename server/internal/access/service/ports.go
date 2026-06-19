package service

import (
	"context"

	accessdb "github.com/findardi/Wadi/server/internal/access/repository/sqlc"
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

	InsertRole(ctx context.Context, arg accessdb.InsertRoleParams) (accessdb.WorkspaceRole, error)

	ExecTx(ctx context.Context, fn func(q *accessdb.Queries) error) error
}

type MailService interface {
	Send(ctx context.Context, to, subject, body string) error
}
