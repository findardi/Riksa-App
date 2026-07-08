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

	CreateDocument(ctx context.Context, arg contentdb.CreateDocumentParams) (contentdb.Document, error)
	CreateDocumentVersion(ctx context.Context, arg contentdb.CreateDocumentVersionParams) (contentdb.DocumentVersion, error)
	SetCurrentVersion(ctx context.Context, arg contentdb.SetCurrentVersionParams) error
	GetNextVersionNo(ctx context.Context, documentID pgtype.UUID) (int32, error)
	GetDocumentByID(ctx context.Context, id pgtype.UUID) (contentdb.Document, error)
	ListDocumentsByFolder(ctx context.Context, folderID pgtype.UUID) ([]contentdb.ListDocumentsByFolderRow, error)
	ListVersionByDocument(ctx context.Context, documentID pgtype.UUID) ([]contentdb.DocumentVersion, error)
	GetVersionByID(ctx context.Context, id pgtype.UUID) (contentdb.DocumentVersion, error)
	GetCurrentVersion(ctx context.Context, id pgtype.UUID) (contentdb.DocumentVersion, error)
	DeleteDocument(ctx context.Context, id pgtype.UUID) error
}
