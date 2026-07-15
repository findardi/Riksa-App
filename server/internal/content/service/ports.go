package service

import (
	"context"

	contentdb "github.com/findardi/Riksa-App/server/internal/content/repository/sqlc"
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

	CreateDocument(ctx context.Context, arg contentdb.CreateDocumentParams) (contentdb.Document, error)
	CreateDocumentVersion(ctx context.Context, arg contentdb.CreateDocumentVersionParams) (contentdb.DocumentVersion, error)
	SetCurrentVersion(ctx context.Context, arg contentdb.SetCurrentVersionParams) error
	GetNextVersionNo(ctx context.Context, documentID pgtype.UUID) (int32, error)
	GetDocumentByID(ctx context.Context, id pgtype.UUID) (contentdb.Document, error)
	ListDocumentsByFolder(ctx context.Context, folderID pgtype.UUID) ([]contentdb.ListDocumentsByFolderRow, error)
	ListVersionByDocument(ctx context.Context, documentID pgtype.UUID) ([]contentdb.DocumentVersion, error)
	GetVersionByID(ctx context.Context, id pgtype.UUID) (contentdb.DocumentVersion, error)
	GetCurrentVersion(ctx context.Context, id pgtype.UUID) (contentdb.DocumentVersion, error)
	SetVersionRendition(ctx context.Context, arg contentdb.SetVersionRenditionParams) error
	DeleteDocument(ctx context.Context, id pgtype.UUID) error
	MoveDocument(ctx context.Context, arg contentdb.MoveDocumentParams) error

	GetAccessLevel(ctx context.Context, arg contentdb.GetAccessLevelParams) (contentdb.AccessLevel, error)
	GetSystemAccessLevelByName(ctx context.Context, name string) (contentdb.AccessLevel, error)
	ListAccessLevels(ctx context.Context, workspaceID pgtype.UUID) ([]contentdb.AccessLevel, error)
	ListFolderAccess(ctx context.Context, arg contentdb.ListFolderAccessParams) ([]contentdb.ListFolderAccessRow, error)
	RemoveFolderAccess(ctx context.Context, arg contentdb.RemoveFolderAccessParams) error
	SetFolderAccess(ctx context.Context, arg contentdb.SetFolderAccessParams) (contentdb.FolderAccess, error)
	ResolveFolderAccess(ctx context.Context, arg contentdb.ResolveFolderAccessParams) (contentdb.ResolveFolderAccessRow, error)
	ListVisibleFolders(ctx context.Context, arg contentdb.ListVisibleFoldersParams) ([]contentdb.ListVisibleFoldersRow, error)

	ExecTx(ctx context.Context, fn func(*contentdb.Queries) error) error
}
