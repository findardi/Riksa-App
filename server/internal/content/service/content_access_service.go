package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/findardi/Riksa-App/server/internal/content/dto"
	contentdb "github.com/findardi/Riksa-App/server/internal/content/repository/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *ContentService) SetFolderAccess(ctx context.Context, req dto.SetFolderAccessRequest) error {
	var wID, fID, gID, lID pgtype.UUID
	if err := wID.Scan(req.WorkspaceID); err != nil {
		return fmt.Errorf("workspace id parse: %w", err)
	}

	if err := fID.Scan(req.FolderID); err != nil {
		return ErrFolderNotFound
	}

	if err := gID.Scan(req.GroupID); err != nil {
		return ErrAccessTargetInvalid
	}

	if err := lID.Scan(req.LevelID); err != nil {
		return ErrAccessTargetInvalid
	}

	_, err := s.repo.SetFolderAccess(ctx, contentdb.SetFolderAccessParams{
		GroupID:     gID,
		WorkspaceID: wID,
		FolderID:    fID,
		LevelID:     lID,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAccessTargetInvalid
	}

	if err != nil {
		return fmt.Errorf("set folder access: %w", err)
	}

	return nil
}

func (s *ContentService) RemoveFolderAccess(ctx context.Context, workspaceID, groupID, folderID string) error {
	var wID, fID, gID pgtype.UUID
	if err := wID.Scan(workspaceID); err != nil {
		return fmt.Errorf("parse workspace id: %w", err)
	}
	if err := fID.Scan(folderID); err != nil {
		return ErrFolderNotFound
	}
	if err := gID.Scan(groupID); err != nil {
		return ErrAccessTargetInvalid
	}

	if err := s.repo.RemoveFolderAccess(ctx, contentdb.RemoveFolderAccessParams{
		FolderID:    fID,
		GroupID:     gID,
		WorkspaceID: wID,
	}); err != nil {
		return fmt.Errorf("remove folder access: %w", err)
	}

	return nil
}

func (s *ContentService) ListAccessLevels(ctx context.Context, workspaceID string) ([]dto.AccessLevelResponse, error) {
	var wID pgtype.UUID
	if err := wID.Scan(workspaceID); err != nil {
		return []dto.AccessLevelResponse{}, fmt.Errorf("workspace id parse: %w", err)
	}

	rows, err := s.repo.ListAccessLevels(ctx, wID)
	if err != nil {
		return []dto.AccessLevelResponse{}, fmt.Errorf("get access level: %w", err)
	}

	access := make([]dto.AccessLevelResponse, 0, len(rows))
	for _, r := range rows {
		access = append(access, dto.AccessLevelResponse{
			ID:           uuidString(r.ID),
			Name:         r.Name,
			IsSystem:     !r.WorkspaceID.Valid,
			CanView:      r.CanView,
			CanDownload:  r.CanDownload,
			CanWatermark: r.CanWatermark,
		})
	}

	return access, nil
}

func (s *ContentService) ListFolderAccess(ctx context.Context, workspaceID, folderID string) ([]dto.FolderAccessResponse, error) {
	var wID, fID pgtype.UUID
	if err := wID.Scan(workspaceID); err != nil {
		return nil, fmt.Errorf("parse workspace id: %w", err)
	}
	if err := fID.Scan(folderID); err != nil {
		return nil, ErrFolderNotFound
	}

	rows, err := s.repo.ListFolderAccess(ctx, contentdb.ListFolderAccessParams{
		FolderID:    fID,
		WorkspaceID: wID,
	})
	if err != nil {
		return nil, fmt.Errorf("list folder access: %w", err)
	}

	res := make([]dto.FolderAccessResponse, 0, len(rows))
	for _, r := range rows {
		res = append(res, dto.FolderAccessResponse{
			FolderID:     uuidString(r.FolderID),
			GroupID:      uuidString(r.GroupID),
			GroupName:    r.GroupName,
			LevelID:      uuidString(r.LevelID),
			LevelName:    r.LevelName,
			CanView:      r.CanView,
			CanDownload:  r.CanDownload,
			CanWatermark: r.CanWatermark,
		})
	}

	return res, nil
}
