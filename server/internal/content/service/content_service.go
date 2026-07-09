package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/findardi/Wadi/server/internal/content/dto"
	contentdb "github.com/findardi/Wadi/server/internal/content/repository/sqlc"
	"github.com/findardi/Wadi/server/internal/platform/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	maxFolderDepth = 256
	uploadURLTTL   = 15 * time.Minute
	downloadURLTTL = 5 * time.Minute
)

var (
	ErrParentCrossWorkspace = errors.New("parent cross workspace")
	ErrParentNotFound       = errors.New("parent not found")
	ErrFolderNameTaken      = errors.New("folder already exists")
	ErrFolderNotFound       = errors.New("folder not found")
	ErrCycle                = errors.New("cannot move folder into its own subtree")
	ErrFolderTreeTooDeep    = errors.New("folder nesting is too deep")
	ErrDocumentNotFound     = errors.New("document not found")
	ErrUploadNotFound       = errors.New("uploaded object not found")
)

type ContentService struct {
	repo  ContentRepository
	store storage.Storage
}

func NewContentService(repo ContentRepository, store storage.Storage) *ContentService {
	return &ContentService{
		repo:  repo,
		store: store,
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

func storageKey(workspaceID, folderID string) string {
	return fmt.Sprintf("%s/%s/%s", workspaceID, folderID, uuid.NewString())
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" && pgErr.ConstraintName == constraint
	}
	return false
}

func (s *ContentService) CreateFolder(ctx context.Context, req dto.CreateFolderRequest) (dto.FolderResponse, error) {

	var wID, pID, cID pgtype.UUID

	if err := cID.Scan(req.CreatedBy); err != nil {
		return dto.FolderResponse{}, fmt.Errorf("user id parse: %w", err)
	}

	if req.ParentID != "" {
		if err := pID.Scan(req.ParentID); err != nil {
			return dto.FolderResponse{}, fmt.Errorf("parent id parse: %w", err)
		}
		pFolder, err := s.repo.GetFolderByID(ctx, pID)
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.FolderResponse{}, ErrParentNotFound
		}

		if err != nil {
			return dto.FolderResponse{}, fmt.Errorf("check parent: %w", err)
		}

		if uuidString(pFolder.WorkspaceID) != req.WorkspaceID {
			return dto.FolderResponse{}, ErrParentCrossWorkspace
		}
	}

	if err := wID.Scan(req.WorkspaceID); err != nil {
		return dto.FolderResponse{}, fmt.Errorf("worspace id parse: %w", err)
	}
	maxPos, err := s.repo.GetMaxPositionInParent(ctx, contentdb.GetMaxPositionInParentParams{
		WorkspaceID: wID,
		ParentID:    pID,
	})

	if err != nil {
		return dto.FolderResponse{}, fmt.Errorf("check max position: %w", err)
	}

	f, err := s.repo.CreateFolder(ctx, contentdb.CreateFolderParams{
		WorkspaceID: wID,
		ParentID:    pID,
		Name:        req.Name,
		Position:    maxPos + 1,
		CreatedBy:   cID,
	})

	if isUniqueViolation(err, "folders_name_root_key") || isUniqueViolation(err, "folders_name_per_parent_key") {
		return dto.FolderResponse{}, ErrFolderNameTaken
	}

	if err != nil {
		return dto.FolderResponse{}, fmt.Errorf("create folder: %w", err)
	}

	return dto.FolderResponse{
		ID:          uuidString(f.ID),
		WorkspaceID: uuidString(f.WorkspaceID),
		ParentID:    uuidString(f.ParentID),
		Name:        f.Name,
		Position:    f.Position,
		CreatedBy:   uuidString(f.CreatedBy),
		CreatedAt:   f.CreatedAt.Time,
		UpdatedAt:   f.UpdatedAt.Time,
	}, nil
}

func (s *ContentService) MoveFolder(ctx context.Context, req dto.MoveFolderRequest) error {

	var fID, pID pgtype.UUID

	if err := fID.Scan(req.FolderID); err != nil {
		return fmt.Errorf("folder id parse: %w", err)
	}

	folder, err := s.repo.GetFolderByID(ctx, fID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFolderNotFound
	}

	if err != nil {
		return fmt.Errorf("get folder: %w", err)
	}

	if req.ParentID != "" {
		if err := pID.Scan(req.ParentID); err != nil {
			return fmt.Errorf("parent id parse: %w", err)
		}

		if pID == fID {
			return ErrCycle
		}

		parent, err := s.repo.GetFolderByID(ctx, pID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrParentNotFound
		}

		if err != nil {
			return fmt.Errorf("check parent: %w", err)
		}

		if parent.WorkspaceID != folder.WorkspaceID {
			return ErrParentCrossWorkspace
		}

		cursor := parent
		for depth := 0; ; depth++ {
			if cursor.ID == fID {
				return ErrCycle
			}

			if !cursor.ParentID.Valid {
				break
			}

			if depth >= maxFolderDepth {
				return ErrFolderTreeTooDeep
			}

			cursor, err = s.repo.GetFolderByID(ctx, cursor.ParentID)
			if err != nil {
				return fmt.Errorf("walk ancestors: %w", err)
			}
		}
	}

	maxPos, err := s.repo.GetMaxPositionInParent(ctx, contentdb.GetMaxPositionInParentParams{
		WorkspaceID: folder.WorkspaceID,
		ParentID:    pID,
	})

	if err != nil {
		return fmt.Errorf("check max position: %w", err)
	}

	err = s.repo.MoveFolder(ctx, contentdb.MoveFolderParams{
		ID:       fID,
		ParentID: pID,
		Position: maxPos + 1,
	})

	if isUniqueViolation(err, "folders_name_root_key") || isUniqueViolation(err, "folders_name_per_parent_key") {
		return ErrFolderNameTaken
	}

	return err
}

func (s *ContentService) GetFoldersTree(ctx context.Context, workspaceID string) ([]dto.FolderTreeNode, error) {

	var wID pgtype.UUID
	if err := wID.Scan(workspaceID); err != nil {
		return nil, fmt.Errorf("workspace id parse: %w", err)
	}

	rows, err := s.repo.GetFoldersByWorkspace(ctx, wID)
	if err != nil {
		return nil, fmt.Errorf("get folders: %w", err)
	}

	childrenOf := make(map[string][]contentdb.Folder)
	for _, f := range rows {
		key := uuidString(f.ParentID)
		childrenOf[key] = append(childrenOf[key], f)
	}

	return buildFolderTree(childrenOf, "", ""), nil
}

func buildFolderTree(childrenOf map[string][]contentdb.Folder, parentKey, prefix string) []dto.FolderTreeNode {
	items := childrenOf[parentKey]
	nodes := make([]dto.FolderTreeNode, 0, len(items))

	for i, f := range items {
		number := prefix + strconv.Itoa(i+1)
		id := uuidString(f.ID)

		nodes = append(nodes, dto.FolderTreeNode{
			ID:       id,
			Name:     f.Name,
			Number:   number,
			Position: f.Position,
			Children: buildFolderTree(childrenOf, id, number+"."),
		})
	}

	return nodes
}

func (s *ContentService) RenameFolder(ctx context.Context, req dto.RenameFolderRequest) (dto.FolderResponse, error) {
	var fID pgtype.UUID
	if err := fID.Scan(req.FolderID); err != nil {
		return dto.FolderResponse{}, fmt.Errorf("folder id parse: %w", err)
	}

	f, err := s.repo.RenameFolder(ctx, contentdb.RenameFolderParams{
		ID:   fID,
		Name: req.Name,
	})

	if isUniqueViolation(err, "folders_name_root_key") || isUniqueViolation(err, "folders_name_per_parent_key") {
		return dto.FolderResponse{}, ErrFolderNameTaken
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return dto.FolderResponse{}, ErrFolderNotFound
	}

	if err != nil {
		return dto.FolderResponse{}, fmt.Errorf("rename folder: %w", err)
	}

	return dto.FolderResponse{
		ID:          uuidString(f.ID),
		WorkspaceID: uuidString(f.WorkspaceID),
		ParentID:    uuidString(f.ParentID),
		Name:        f.Name,
		Position:    f.Position,
		CreatedBy:   uuidString(f.CreatedBy),
		CreatedAt:   f.CreatedAt.Time,
		UpdatedAt:   f.UpdatedAt.Time,
	}, nil
}

func (s *ContentService) DeleteFolder(ctx context.Context, folderID string) error {
	var fID pgtype.UUID
	if err := fID.Scan(folderID); err != nil {
		return fmt.Errorf("folder id parse: %w", err)
	}

	if _, err := s.repo.GetFolderByID(ctx, fID); errors.Is(err, pgx.ErrNoRows) {
		return ErrFolderNotFound
	} else if err != nil {
		return fmt.Errorf("get folder: %w", err)
	}

	if err := s.repo.DeleteFolder(ctx, fID); err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	return nil
}
