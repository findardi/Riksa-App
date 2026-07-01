package service

import (
	"context"

	contentdb "github.com/findardi/Wadi/server/internal/content/repository/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type ContentRepository interface {
	CreateFolder(ctx context.Context, arg contentdb.CreateFolderParams) (contentdb.Folder, error)
	DeleteFolder(ctx context.Context, id pgtype.UUID) error
	GetFolderByID(ctx context.Context, id pgtype.UUID) (contentdb.Folder, error)
	GetFoldersByWorkspace(ctx context.Context, workspaceID pgtype.UUID) ([]contentdb.Folder, error)
	GetMaxPositionInParent(ctx context.Context, arg contentdb.GetMaxPositionInParentParams) (int32, error)
	MoveFolder(ctx context.Context, arg contentdb.MoveFolderParams) error
	RenameFolder(ctx context.Context, arg contentdb.RenameFolderParams) (contentdb.Folder, error)
}
