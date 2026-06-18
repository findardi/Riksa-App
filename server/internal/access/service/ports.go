package service

import (
	"context"

	accessdb "github.com/findardi/Wadi/server/internal/access/repository/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessRepository interface {
	AddMember(ctx context.Context, arg accessdb.AddMemberParams) (accessdb.WorkspaceMember, error)

	DeleteRole(ctx context.Context, id pgtype.UUID) error

	EditRole(ctx context.Context, arg accessdb.EditRoleParams) (accessdb.WorkspaceRole, error)

	GetRole(ctx context.Context, id pgtype.UUID) (accessdb.WorkspaceRole, error)
	GetRoles(ctx context.Context, workspaceID pgtype.UUID) ([]accessdb.WorkspaceRole, error)
	InsertRole(ctx context.Context, arg accessdb.InsertRoleParams) (accessdb.WorkspaceRole, error)

	ExecTx(ctx context.Context, fn func(q *accessdb.Queries) error) error
}

type MailService interface {
	Send(ctx context.Context, to, subject, body string) error
}
